package internal

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
)

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	p, err := app.Forum.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &TemplateData{
		IsAuth: app.CookieGet(r),
		Posts:  p,
	})
}

func (app *App) Filter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filter" {
		app.notFound(w)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.badRequest(w)
		return
	}

	values, err := url.ParseQuery(string(bytes))
	if err != nil {
		app.badRequest(w)
		return
	}
	if _, ok := values["filter"]; !ok {
		app.badRequest(w)
		return
	}
	tag := ""
	for i, v := range values {
		switch i {
		case "filter":
			tag = v[0]
		}
	}

	p, err := app.Forum.GetPostsWithTag(tag)
	if err != nil {
		app.serverError(w, err)
	}
	app.render(w, r, "home.page.html", &TemplateData{
		IsAuth: app.CookieGet(r),
		Posts:  p,
	})
}

func (app *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createpost" {
		app.notFound(w)
		return
	}
	if !app.CookieGet(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "createpost.page.html", &TemplateData{IsAuth: app.CookieGet(r)})
	case http.MethodPost:
		post := models.Post{
			Created: time.Now().Format("02 Jan 2006 15:04:05"),
		}
		t := models.Tag{}

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.badRequest(w)
			return
		}

		values, err := url.ParseQuery(string(bytes))
		if err != nil {
			app.badRequest(w)
			return
		}
		if _, ok := values["title"]; !ok {
			app.badRequest(w)
			return
		}
		if _, ok := values["content"]; !ok {
			app.badRequest(w)
			return
		}

		for i, v := range values {
			switch i {
			case "title":
				if !checkInput(v[0]) {
					app.badRequest(w)
					return
				}
				post.Title = v[0]
			case "content":
				if !checkInput(v[0]) {
					app.badRequest(w)
					return
				}
				post.Content = v[0]
			case "tag1":
				t.Hashtags = append(t.Hashtags, v[0])

			case "tag2":
				t.Hashtags = append(t.Hashtags, v[0])
			case "tag3":
				t.Hashtags = append(t.Hashtags, v[0])
			case "tag4":
				t.Hashtags = append(t.Hashtags, v[0])
			}
		}

		if !checkTag(t.Hashtags) {
			app.badRequest(w)
			return
		}

		c, _ := r.Cookie(CookieName)

		u, err := app.Forum.GetUserByUuid(c.Value)
		if err != nil {
			app.ErrorLog.Fatal(err)
			return
		}
		post.Author = u.Username

		postid, err := app.Forum.InsertPost(&post)
		if err != nil {
			app.serverError(w, err)
			return
		}

		post.PostID = int(postid)
		for i := 0; i < len(t.Hashtags); i++ {
			err = app.Forum.InsertTag(post.PostID, t.Hashtags[i])
			if err != nil {
				app.serverError(w, err)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func checkInput(s string) bool {
	if len(RemoveSpace(s)) == 0 {
		return false
	}
	return true
}

func checkTag(s []string) bool {
	arr := []string{"dogs", "cats", "fishes", "birds"}

	for _, v1 := range s {
		ok := false
		for _, v2 := range arr {
			if v1 == v2 {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i != j {
				if arr[i] == arr[j] {
					return false
				}
			}
		}
	}
	return true
}

func (app *App) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		app.notFound(w)
		return
	}
	if app.CookieGet(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "signup.page.html", &TemplateData{})
	case http.MethodPost:
		user := models.User{}

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.badRequest(w)
			return
		}

		values, err := url.ParseQuery(string(bytes))
		if err != nil {
			app.badRequest(w)
			return
		}
		if _, ok := values["username"]; !ok {
			app.badRequest(w)
			return
		}
		if _, ok := values["email"]; !ok {
			app.badRequest(w)
			return
		}
		if _, ok := values["password"]; !ok {
			app.badRequest(w)
			return
		}
		for i, v := range values {
			switch i {
			case "username":
				user.Username = v[0]
			case "email":
				user.Email = v[0]
			case "password":
				user.Password = v[0]
			}
		}
		if !isEmailValid(user.Email) {
			app.ErrorLog.Println("incorrect email")

			app.render(w, r, "signup.page.html", &TemplateData{
				Err: Error{
					IsError: true,
					Msg:     "incorrect email",
				},
			})
			return
		}
		if !isUsernameValid(user.Username) {
			app.ErrorLog.Println("incorrect username")

			app.render(w, r, "signup.page.html", &TemplateData{
				Err: Error{
					IsError: true,
					Msg:     "incorrect username",
				},
			})
			return
		}
		if len(user.Password) < 5 {
			app.ErrorLog.Println("incorrect pass")

			app.render(w, r, "signup.page.html", &TemplateData{
				Err: Error{
					IsError: true,
					Msg:     "minimum 5 sized password",
				},
			})
			return
		}
		_, err = app.Forum.GetUserByName(user.Username)
		if err == nil {
			app.ErrorLog.Println("username taken")

			app.render(w, r, "signup.page.html", &TemplateData{
				Err: Error{
					IsError: true,
					Msg:     "username taken",
				},
			})
			return
		}

		_, err = app.Forum.GetUserByEmail(user.Email)
		if err == nil {
			app.ErrorLog.Println("email taken")
			app.badRequest(w)
			app.render(w, r, "signup.page.html", &TemplateData{
				Err: Error{
					IsError: true,
					Msg:     "email taken",
				},
			})
			return
		}
		user.Password, _ = HashPassword(user.Password)

		userID, err := app.Forum.InsertUser(user)
		if err != nil {
			app.serverError(w, err)
			return
		}
		user.UserID = int(userID)
		app.CookieSet(w, r, user.UserID)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		app.methodNotAllowed(w)
	}
}

func (app *App) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		app.notFound(w)
		return
	}
	if app.CookieGet(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "signin.page.html", &TemplateData{})
	case http.MethodPost:
		user := models.User{}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			app.badRequest(w)
			return
		}

		values, err := url.ParseQuery(string(bytes))
		if err != nil {
			app.badRequest(w)
			return
		}
		if _, ok := values["username"]; !ok {
			app.badRequest(w)
			return
		}
		if _, ok := values["password"]; !ok {
			app.badRequest(w)
			return
		}

		for i, v := range values {
			switch i {
			case "username":
				user.Username = v[0]
			case "password":
				user.Password = v[0]
			case "login":
				if v[0] != "login" {
					app.badRequest(w)
					return
				}
			default:
				app.badRequest(w)
				return
			}
		}
		u, err := app.Forum.GetUserByName(user.Username)
		if err != nil || !CheckPasswordHash(user.Password, u.Password) {
			app.ErrorLog.Println(err)
			app.render(w, r, "signin.page.html", &TemplateData{
				Err: Error{
					IsError: true,
					Msg:     "wrong username/pass",
				},
			})
			return
		}

		app.Forum.DelCookie(u.UserID)

		app.CookieSet(w, r, u.UserID)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		app.methodNotAllowed(w)
	}
}

func (app *App) SignOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signout" {
		app.notFound(w)
		return
	}
	if !app.CookieGet(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   CookieName,
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) Profile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile" {
		app.notFound(w)
		return
	}
	if !app.CookieGet(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	switch r.Method {
	case http.MethodGet:
		c, _ := r.Cookie(CookieName)

		u, err := app.Forum.GetUserByUuid(c.Value)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}
		p, err := app.Forum.GetMyPosts(u.Username)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}
		lp, err := app.Forum.GetLikedPosts(u.UserID)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}
		app.render(w, r, "profile.page.html", &TemplateData{
			IsAuth:     app.CookieGet(r),
			User:       u,
			Posts:      p,
			LikedPosts: lp,
		})
	default:
		app.methodNotAllowed(w)
	}
}

func (app *App) ShowPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/showpost" {
		app.notFound(w)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.methodNotAllowed(w)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.badRequest(w)
		return
	}

	p, err := app.Forum.GetPostById(id)
	if err != nil {
		app.notFound(w)
		// app.serverError(w, err)
		return
	}

	app.render(w, r, "showpost.page.html", &TemplateData{
		IsAuth: app.CookieGet(r),
		Post:   p,
	})
}

func (app *App) RatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.methodNotAllowed(w)
		return
	}

	if !app.CookieGet(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	p, err := app.Forum.GetPostById(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var symbol string
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.badRequest(w)
		return
	}

	values, err := url.ParseQuery(string(bytes))
	if err != nil {
		app.badRequest(w)
		return
	}
	if _, ok := values["like"]; !ok {
		app.badRequest(w)
		return
	}
	for i, v := range values {
		switch i {
		case "like":
			symbol = v[0]
		}
	}

	c, err := r.Cookie(CookieName)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	u, err := app.Forum.GetUserByUuid(c.Value)
	if err != nil {
		app.ErrorLog.Fatal(err)
		return
	}
	sym, err := app.Forum.GetRaiting(u.UserID, p.PostID)
	if err != nil {
		if err := app.Forum.InsertRaiting(u.UserID, p.PostID, symbol); err != nil {
			app.serverError(w, err)
			return
		}
	} else {
		if err := app.Forum.DelRaiting(u.UserID, p.PostID); err != nil {
			app.serverError(w, err)
			return
		}
		if sym != symbol {
			if err := app.Forum.InsertRaiting(u.UserID, p.PostID, symbol); err != nil {
				app.serverError(w, err)
				return
			}
		}
	}

	http.Redirect(w, r, "/showpost?id="+strconv.Itoa(id), http.StatusSeeOther)
}

func (app *App) rateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.methodNotAllowed(w)
		return
	}

	if !app.CookieGet(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	p, err := app.Forum.GetPostById(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var symbol string
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.badRequest(w)
		return
	}

	values, err := url.ParseQuery(string(bytes))
	if err != nil {
		app.badRequest(w)
		return
	}
	if _, ok := values["like"]; !ok {
		app.badRequest(w)
		return
	}
	for i, v := range values {
		switch i {
		case "like":
			symbol = v[0]
		}
	}
	symbols := strings.Split(symbol, "and")
	if len(symbols) != 2 {
		app.badRequest(w)
		return
	}
	num, err := strconv.Atoi(symbols[0])
	if err != nil {
		app.badRequest(w)
		return
	}

	c, err := r.Cookie(CookieName)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	u, err := app.Forum.GetUserByUuid(c.Value)
	if err != nil {
		app.ErrorLog.Fatal(err)
		return
	}

	sym, err := app.Forum.GetRateComment(u.UserID, p.PostID, num)
	if err != nil {
		if err := app.Forum.InsertRateComment(u.UserID, p.PostID, num, symbols[1]); err != nil {
			app.serverError(w, err)
			return
		}
	} else {
		if err := app.Forum.DelRateComment(u.UserID, num); err != nil {
			app.serverError(w, err)
			return
		}
		if sym != symbols[1] {
			if err := app.Forum.InsertRateComment(u.UserID, p.PostID, num, symbols[1]); err != nil {
				app.serverError(w, err)
				return
			}
		}
	}
	http.Redirect(w, r, "/showpost?id="+strconv.Itoa(id), http.StatusSeeOther)
}

func (app *App) CommentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.methodNotAllowed(w)
		return
	}
	if !app.CookieGet(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	p, err := app.Forum.GetPostById(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	c, err := r.Cookie(CookieName)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	u, err := app.Forum.GetUserByUuid(c.Value)
	if err != nil {
		app.ErrorLog.Fatal(err)
		return
	}

	var text string
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.badRequest(w)
		return
	}

	values, err := url.ParseQuery(string(bytes))
	if err != nil {
		app.badRequest(w)
		return
	}

	if _, ok := values["input"]; !ok {
		app.badRequest(w)
		return
	}

	for i, v := range values {
		switch i {
		case "input":
			text = v[0]
		}
	}
	if !checkInput(text) {
		app.badRequest(w)
		return
	}

	err = app.Forum.InsertComment(u.UserID, u.Username, p.PostID, text)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/showpost?id="+strconv.Itoa(id), http.StatusSeeOther)
}
