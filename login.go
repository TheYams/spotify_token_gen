package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		Scopes:       []string{"playlist-read-private", "playlist-read-collaborative"},
		RedirectURL:  os.Getenv("SPOTIFY_REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}

	token := fetchToken(conf)
	saveToken(token)
	_ = conf.Client(oauth2.NoContext, token)
}

func saveToken(token *oauth2.Token) {
	b, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("tokens.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(b))
}

func fetchToken(conf *oauth2.Config) *oauth2.Token {
	state := "ideallygenerated"

	showDialogOption := oauth2.SetAuthURLParam("show_dialog", "false")

	url := conf.AuthCodeURL(state, showDialogOption)
	log.Print("Visit the URL for the auth dialog: ", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Auth Token is: ", token.AccessToken)
	log.Println("Refresh Token is:", token.RefreshToken)

	return token
}

func loadToken() (token *oauth2.Token, err error) {
	b, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(b, &token)
	if err != nil {
		return token, err
	}

	return token, nil
}
