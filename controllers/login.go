package controllers

import (
	"GoogleLogin/models"
	"encoding/json"
	"fmt"
	"github.com/kataras/go-sessions"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"time"
)

const htmlIndex = `<html><body><a href="/GoogleLogin">Log in with Google</a></body></html>`

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL: "http://localhost:9000/GoogleCallback",
		ClientID:    "",
		ClientSecret: "	",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	oauthStateString = "random"
)

func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	URL := "static/redirect/path"
	fmt.Println(URL)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	UserObj := make(map[string]interface{})
	json.Unmarshal(contents, &UserObj)
	GoogleUserTable, err := models.GetGoogleUserTableByAuthId(cast.ToString(UserObj["id"]))

	if err != nil {
		newGoogleUserTable := models.GoogleUserTable{}
		newGoogleUserTable.AuthId = cast.ToString(UserObj["id"])
		newGoogleUserTable.VerifiedEmail = cast.ToInt8(UserObj["verified_email"])
		newGoogleUserTable.Email = cast.ToString(UserObj["email"])
		newGoogleUserTable.Hd = cast.ToString(UserObj["hd"])
		newGoogleUserTable.Name = cast.ToString(UserObj["name"])
		newGoogleUserTable.Picture = cast.ToString(UserObj["picture"])
		newGoogleUserTable.Role = "user"

		models.AddGoogleUserTable(&newGoogleUserTable)
		sess := sessions.Start(w, r)
		sess.Set(cast.ToString(UserObj["id"]), newGoogleUserTable)
		fmt.Println("user not exist ")
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "frontend", Value: cast.ToString(UserObj["id"]), Expires: expiration}
		http.SetCookie(w, &cookie)

		fmt.Println(sess.Get(cast.ToString(UserObj["id"])))
		//http.Redirect(w, r, URL, http.StatusTemporaryRedirect)

	} else {
		sess := sessions.Start(w, r)
		sess.Set(cast.ToString(UserObj["id"]), GoogleUserTable)

		fmt.Println("user exist already")
		expiration := time.Now().Add(365 * 24 * time.Hour)

		cookie := http.Cookie{Name: "frontend", Value: cast.ToString(UserObj["id"]), Expires: expiration}
		http.SetCookie(w, &cookie)
		fmt.Println(sess.Get(cast.ToString(UserObj["id"])))
		//http.Redirect(w, r, URL, http.StatusTemporaryRedirect)

	}
	returnJson, _ := json.Marshal(UserObj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnJson)
}
