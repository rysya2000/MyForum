package internal

import (
	"net/http"
)

func (app *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Home)

	mux.HandleFunc("/createpost", app.CreatePost)
	mux.HandleFunc("/showpost", app.ShowPost)
	mux.HandleFunc("/showpost/rate", app.RatePost)
	mux.HandleFunc("/showpost/comment", app.CommentPost)
	mux.HandleFunc("/showpost/rateComment", app.rateComment)

	mux.HandleFunc("/signup", app.SignUp)
	mux.HandleFunc("/signin", app.SignIn)
	mux.HandleFunc("/signin/github/", app.githubLoginHandler)
	mux.HandleFunc("/signin/github/callback", app.githubCallbackHandler)
	mux.HandleFunc("/signin/google/", app.googleLoginHandler)
	mux.HandleFunc("/signin/google/callback", app.googleCallbackHandler)
	mux.HandleFunc("/signout", app.SignOut)

	mux.HandleFunc("/profile", app.Profile)
	mux.HandleFunc("/filter", app.Filter)

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))))
	return mux
}
