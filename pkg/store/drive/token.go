package drive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err

	}
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	return tok, nil

}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveToken(file string, token *oauth2.Token) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func GetToken(config *oauth2.Config, cacheFile string) (*oauth2.Token, error) {
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err := saveToken(cacheFile, tok); err != nil {
			return nil, err
		}
	}
	return tok, nil
}

func GetClient(ctx context.Context, scope, credFile, cacheFile string) (*http.Client, error) {
	b, err := ioutil.ReadFile(credFile)
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		return nil, err
	}
	token, err := GetToken(config, cacheFile)
	if err != nil {
		return nil, err
	}
	return config.Client(ctx, token), nil
}
