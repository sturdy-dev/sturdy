package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"client/cmd/sturdy/config"
	"client/cmd/sturdy/legal"
	"client/cmd/sturdy/version"
)

func printHelpAndExit() {
	fmt.Println("This is the Sturdy client")
	fmt.Println("")
	fmt.Println("No subcommand provided")
	fmt.Println("")
	fmt.Println("Available commands:")
	fmt.Println("  start    Start Sturdy connections for all connected codebases")
	fmt.Println("  stop     Stop all connections and stop the daemon")
	fmt.Println("  restart  Restart and re-configure all connections")
	fmt.Println("  status   Get the current status of each codebase")
	fmt.Println("  auth     Authenticate yourself with Sturdy")
	fmt.Println("  init     Configure a new codebase to be used from this computer")
	fmt.Println("  import   Import a Git repository to Sturdy")
	fmt.Println("  version  Display Sturdy version information")
	fmt.Println("  legal    Display legal credits")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printHelpAndExit()
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Could not find user home dir")
		os.Exit(1)
		return
	}

	// len is 2 if there is only a subcommand provided, and no additional flags
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	fs := flag.FlagSet{}
	configPath := fs.String("config", path.Join(home, ".sturdy"), "Path to your Sturdy configuration file")
	err = fs.Parse(args)
	if err != nil {
		log.Println("Failed to parse flags", err)
		os.Exit(1)
		return
	}

	// Remaining arguments after flags have been parsed
	args = fs.Args()

	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
		return
	}

	exitIfErr := func(err error) {
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	switch os.Args[1] {
	case "auth":
		// It's important to not attempt to require auth, or renew auth _before_ calling auth()
		auth(conf, *configPath)
	case "status":
		status(conf)
	case "init":
		apiClient, err := requireAuth(conf, *configPath)
		exitIfErr(err)
		initSmart(conf, *configPath, args, apiClient)
	case "start":
		apiClient, err := requireAuth(conf, *configPath)
		exitIfErr(err)
		startMutagen(*configPath, conf, apiClient)
	case "stop":
		stopMutagen(conf)
	case "restart":
		stopMutagen(conf)
		apiClient, err := requireAuth(conf, *configPath)
		exitIfErr(err)
		startMutagen(*configPath, conf, apiClient)
	case "legal":
		fmt.Println(legal.LegalNotice)
	case "version":
		version.VersionCMD()
	case "import":
		apiClient, err := requireAuth(conf, *configPath)
		exitIfErr(err)
		importCodebase(conf, args, apiClient)
	default:
		printHelpAndExit()
	}
}
