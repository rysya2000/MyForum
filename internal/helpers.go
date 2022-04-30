package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
)

func (app *App) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *App) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *App) methodNotAllowed(w http.ResponseWriter) {
	app.clientError(w, http.StatusMethodNotAllowed)
}

func (app *App) badRequest(w http.ResponseWriter) {
	app.clientError(w, http.StatusBadRequest)
}

func (app *App) render(w http.ResponseWriter, r *http.Request, name string, td interface{}) {
	err := app.TemplateCache.ExecuteTemplate(w, name, td)
	if err != nil {
		app.serverError(w, err)
	}
}

func ParseEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, v := range lines {
		t := strings.Split(v, "=")
		if len(t) != 2 {
			return errors.New("wrong env")
		}
		os.Setenv(t[0], t[1])
		os.Getenv(t[0])
	}

	return nil
}

func getGithubAccesToken(code string) string {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECURE")
	if clientID == "" || clientSecret == "" {
		log.Fatal("getenv failed")
	}

	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}

	requestJSON, _ := json.Marshal(requestBodyMap)

	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		log.Fatal("Request creation failed")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}
	respbody, _ := ioutil.ReadAll(resp.Body)

	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)
	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	type Data struct {
		Username string `json:"login"`
		Email    string `json:"email"`
	}
	var data Data

	json.Unmarshal(respbody, &data)

	u := models.User{
		Username: data.Username,
	}
	return u.Username
}

func getGithubEmail(accessToken string) string {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	type Data struct {
		Username string `json:"given_name"`
		Email    string `json:"email"`
	}
	var data []Data

	json.Unmarshal(respbody, &data)

	u := models.User{
		Email: data[0].Email,
	}
	return u.Email
}

func getGoogleAccessToken(code string) string {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECURE")
	if clientID == "" || clientSecret == "" {
		log.Panic("getenv failed")
	}

	u := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {"http://localhost:8000/signin/google/callback"},
	}

	req, err := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(u.Encode()),
	)
	if err != nil {
		log.Panic("Request creation failed")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	type googleAccesTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	var ghresp googleAccesTokenResponse
	err = json.Unmarshal(respbody, &ghresp)
	if err != nil {
		log.Println(err)
	}
	if ghresp.AccessToken == "" {
		log.Printf("empty access token")
	}
	return ghresp.AccessToken
}

func getGoogleData(accessToken string) models.User {
	req, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v3/userinfo?access_token="+accessToken,
		nil,
	)
	if err != nil {
		log.Panic("NewReq failed")
	}

	auth := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("ReadAll failed")
	}

	type Data struct {
		Username string `json:"given_name"`
		Email    string `json:"email"`
	}
	var data Data

	json.Unmarshal(body, &data)

	u := models.User{
		Username: data.Username,
		Email:    data.Email,
	}

	return u
}

func generateRandomStrongPassword() string {
	// random set of: abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789!@#$&*_
	charSet := "RW!Vx&6fMHyPYSB1da*4LTm3ki@5c2ptgDzZ9Gq8w7Ke$XNE#s_jvrJuQnFCAUbh"
	pass := randomStringGenerator(charSet, 14)

	return pass
}

func randomStringGenerator(charSet string, codeLength int32) string {
	code := ""
	charSetLength := int32(len(charSet))
	for i := int32(0); i < codeLength; i++ {
		index := randomNumber(0, charSetLength)
		code += string(charSet[index])
	}

	return code
}

func randomNumber(min, max int32) int32 {
	rand.Seed(time.Now().UnixNano())
	return min + int32(rand.Intn(int(max-min)))
}
