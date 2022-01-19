package main

import (
	"fmt"
	"log"

	"client/cmd/sturdy/config"
	"client/pkg/mutagen"
)

func status(conf *config.Config) {
	mutagenSessions, err := mutagen.Status()
	if err != nil {
		log.Fatal(err)
	}

	sessionsByName := make(map[string]mutagen.SessionStatus)
	for _, s := range mutagenSessions {
		sessionsByName[s.Session.Name] = s
	}

	for _, view := range conf.Views {
		name := viewMutagenName(view)

		session, ok := sessionsByName[name]
		if !ok {
			fmt.Printf("❓ %s - no session found\n", view.Path)
			continue
		}

		if session.AlphaConnected && session.BetaConnected {
			fmt.Printf("\U0001F7E2 %s is connected\n", view.Path)
			continue
		}

		if session.Session.Paused {
			fmt.Printf("⏸️  %s is paused\n", view.Path)
			continue
		}

		fmt.Printf("❓ %s\tunknown status\n", view.Path)
	}
}
