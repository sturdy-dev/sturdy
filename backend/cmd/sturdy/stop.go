package main

import (
	"fmt"
	"log"
	"mash/cmd/sturdy/config"
	"mash/cmd/sturdy/pkg/mutagen"
)

func stopMutagen(conf *config.Config) {
	// Check if there are any stray mutagen connections
	mutagenSessions, err := mutagen.Status()
	if err != nil {
		fmt.Printf("⚠️  Failed to check status of connections\n")
	}
	sessionsByName := make(map[string]mutagen.SessionStatus)
	for _, s := range mutagenSessions {
		sessionsByName[s.Session.Name] = s
	}

	// Pause mutagen for each view
	for _, view := range conf.Views {
		name := viewMutagenName(view)

		// If there's no mutagen session, don't stop it
		session, ok := sessionsByName[name]
		if !ok || session.Session.Paused {
			fmt.Printf("⏸️  Stopped %s (already stopped)\n", view.Path)
			continue
		}

		err := terminateMutagenForView(view)
		if err != nil {
			fmt.Printf("⚠️  Could not stop %s\n", view.Path)
			log.Println(err)
			continue
		}
		fmt.Printf("⏸️  Stopped %s\n", view.Path)
	}

mutagenStatuses:
	for _, ms := range mutagenSessions {
		// Connection not managed by Sturdy?
		if !isSturdySession(ms) {
			continue
		}
		if ms.Session.Paused {
			continue
		}

		// Check if we have a entry for this "sync" in our config
		for _, view := range conf.Views {
			name := viewMutagenName(view)
			if ms.Session.Name == name {
				continue mutagenStatuses
			}
		}

		// No matched name
		fmt.Printf("⚠️  Found directory with missing configuration, pausing it now (%s)\n", ms.Session.Alpha.Path)
		err = mutagenPauseSyncByName(ms.Session.Identifier)
		if err != nil {
			fmt.Printf("⚠️  Could not stop %s\n", ms.Session.Name)
			log.Println(err)
		}
	}

	// Stop the daemon
	_, err = mutagen.RunMutagenCommandWithRestart("daemon", "stop")
	if err != nil {
		fmt.Printf("⚠️  Could not stop daemon\n")
		log.Println(err)
		return
	}

	fmt.Printf("⏸️  Sturdy has been stopped\n")
	return
}

func isSturdySession(status mutagen.SessionStatus) bool {
	if status.Session.Beta.Host == "sync.getsturdy.com" {
		return true
	}
	if _, ok := status.Session.Labels["sturdy"]; ok {
		return true
	}
	return false
}

func mutagenPauseSync(view config.ViewConfig) error {
	return mutagenPauseSyncByName(viewMutagenName(view))
}

func mutagenPauseSyncByName(name string) error {
	_, err := mutagen.RunMutagenCommandWithRestart("sync", "pause", name)
	if err != nil {
		return fmt.Errorf("failed to pause: %w", err)
	}
	return nil
}
