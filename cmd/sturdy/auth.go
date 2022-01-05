package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"mash/cmd/sturdy/pkg/api"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"mash/cmd/sturdy/config"
)

func auth(conf *config.Config, configPath string) {
	url := "https://getsturdy.com/install/token"
	fmt.Printf("👉 Open this page in your browser: %s\n", url)
	fmt.Print("🔑 Paste the code that the browser gave you, and press enter")

	code, err := readUntilValidTokenInput(conf, os.Stdin, checkToken)
	if errors.Is(err, io.EOF) {
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Failed read the auth code. Please try running this command again.")
		fmt.Println(err)
		os.Exit(1)
	}

	err = config.SetAuth(configPath, code)
	if err != nil {
		fmt.Println("Failed to update config")
		fmt.Println(err)
		return
	}
	conf.Auth = code

	// Create a new API client
	apiClient := api.NewHttpApiClient(conf)
	user, err := apiClient.GetUser()
	if err != nil {
		fmt.Println("Something went wrong")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Success! You're connected as %s.\n", user.Name)
	fmt.Printf("The configuration has been saved to %s\n", configPath)
}

func readUntilValidTokenInput(conf *config.Config, termReadWriter io.ReadWriter, validateFunc validateTokenFunc) (string, error) {
	fmt.Println()

	for attempt := 0; attempt < 30; attempt++ {

		reader := bufio.NewReader(termReadWriter)
		fmt.Print(" > ")
		codeBytes, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			return "", err
		}
		if err != nil {
			fmt.Println(err)
			fmt.Println("❌ Something went wrong reading the input. Please try pasting again.")
			continue
		}

		input := strings.TrimSpace(codeBytes)
		err = validateFunc(conf, input)
		if err != nil {
			fmt.Println("❌ Invalid token. Please try pasting it again.")
			continue
		}

		return input, nil
	}

	fmt.Println("\r")

	return "", errOutOfTries
}

var errOutOfTries = fmt.Errorf("❌ Maximum attempts reached, aborting!")

type validateTokenFunc func(conf *config.Config, checkToken string) error

func checkToken(conf *config.Config, checkToken string) error {
	copyConf := *conf
	copyConf.Auth = checkToken
	apiClient := api.NewHttpApiClient(&copyConf)
	_, err := apiClient.GetUser()
	if err != nil {
		return err
	}
	return nil
}

func renewAuth(conf *config.Config, configPath string, api api.SturdyAPI) error {
	// Not authed, don't do anything
	if len(conf.Auth) == 0 {
		return nil
	}

	// This does _not_ validate the token. It simply extracts the expiration date to check if we're eligible for a token renewal
	token, _, err := new(jwt.Parser).ParseUnverified(conf.Auth, jwt.MapClaims{})
	if err != nil {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	exp, ok := claims["exp"]
	if !ok {
		return nil
	}

	expFloat, ok := exp.(float64)
	if !ok {
		return nil
	}

	expT := time.Unix(int64(expFloat), 0)

	// If expires in more than 25 days
	if !expT.Before(time.Now().Add(time.Hour * 24 * 25)) {
		return nil
	}

	res, err := api.RenewAuth()
	if err != nil {
		return fmt.Errorf("failed to renew authentication credentials, if this keeps happening, run 'sturdy auth' and follow the instructions: %w", err)
	}

	if !res.HasNew || len(res.Token) < 10 {
		return nil
	}

	// Updated token!
	conf.Auth = res.Token

	err = config.WriteConfig(configPath, conf)
	if err != nil {
		return fmt.Errorf("failed to renew authentication credentials, if this keeps happening, run 'sturdy auth' and follow the instructions: %w", err)
	}

	return nil
}

func requireAuth(conf *config.Config, configPath string) (*api.HttpApiClient, error) {
	// New authentication
	if conf.Auth == "" {
		auth(conf, configPath)
		return api.NewHttpApiClient(conf), nil
	}

	// Check if we need to renew the authentication
	apiClient := api.NewHttpApiClient(conf)
	err := renewAuth(conf, configPath, apiClient)
	if err != nil {
		return nil, err
	}

	return api.NewHttpApiClient(conf), nil
}
