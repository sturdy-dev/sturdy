package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(tokenJsonPath string, config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokenJsonPath)
	if err != nil {
		// TODO: Token creation path
		// tok = getTokenFromWeb(config)
		// saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func GetService(credentialsJsonPath, tokenJsonPath string) (*gmail.Service, error) {
	b, err := ioutil.ReadFile(credentialsJsonPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	client := getClient(tokenJsonPath, config)

	srv, err := gmail.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client: %w", err)
	}

	return srv, nil
}
