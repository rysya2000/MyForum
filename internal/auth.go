package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
)

func (app *App) githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	if clientID == "" {
		app.ErrorLog.Fatal("Github client id not defined in .env file")
	}

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		clientID,
		"http://localhost:8000/signin/github/callback",
	)

	http.Redirect(w, r, redirectURL, 301)
}

func (app *App) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := getGithubAccesToken(code)

	githubLogin := getGithubData(githubAccessToken)

	githubEmail := getGithubEmail(githubAccessToken)

	if githubLogin == "" || githubEmail == "" {
		app.ErrorLog.Fatal("getting empty login or email ")
	}

	u := models.User{
		Username: githubLogin,
		Email:    githubEmail,
		Password: generateRandomStrongPassword(),
	}

	app.InfoLog.Printf("githubData: %v %v %v\n", u.Username, u.Email, u.Password)

	app.SetUser(w, r, u)
}

func (app *App) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		log.Fatal("Google client id not defined in .env file")
	}

	redirectURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=%s",
		clientID,
		"http://localhost:8000/signin/google/callback",
		"profile email",
		"code",
	)

	http.Redirect(w, r, redirectURL, 301)
}

func (app *App) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	googleAccesToken := getGoogleAccessToken(code)

	googleData := getGoogleData(googleAccesToken)

	if googleData.Username == "" || googleData.Email == "" {
		app.ErrorLog.Fatal("getting empty login or email ")
	}

	u := models.User{
		Username: googleData.Username,
		Email:    googleData.Email,
		Password: generateRandomStrongPassword(),
	}

	app.InfoLog.Printf("User: %v\n", u)

	app.SetUser(w, r, u)
}

func (app *App) SetUser(w http.ResponseWriter, r *http.Request, u models.User) {
	user, err := app.Forum.GetUserByEmail(u.Email)
	if err == sql.ErrNoRows {
		userID, err := app.Forum.InsertUser(u)
		if err != nil {
			app.serverError(w, err)
			return
		}
		u.UserID = int(userID)
		app.CookieSet(w, r, u.UserID)

	} else if err == nil {
		app.Forum.DelCookie(user.UserID)
		app.CookieSet(w, r, user.UserID)
	}
	fmt.Println("HERERERERER")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
