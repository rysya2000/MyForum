package internal

import (
	"fmt"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

const CookieName string = "forum"

func (app *App) CookieSet(w http.ResponseWriter, r *http.Request, nameid int) {
	u := uuid.NewV4()

	app.Forum.InsertCookie(nameid, u.String())

	http.SetCookie(w, &http.Cookie{
		Name:   CookieName,
		Value:  u.String(),
		MaxAge: 1800,
		Path:   "/",
	})
}

func (app *App) CookieGet(r *http.Request) bool {
	c, err := r.Cookie(CookieName)
	if err != nil {
		log.Println("cookieGet: ", err)
		return false
	}
	u, err := app.Forum.GetUserByUuid(c.Value)
	if err != nil {
		fmt.Println("Here2")
		return false
	}
	err = app.Forum.IsCookieInDB(u.UserID)
	if err != nil {
		fmt.Println("Here3")
		return false
	}

	return true
}
