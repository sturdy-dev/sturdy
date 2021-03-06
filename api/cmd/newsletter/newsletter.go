package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/keighl/postmark"

	"getsturdy.com/api/pkg/newsletter"
)

func main() {
	serverToken := flag.String("server-token", "", "Postmark server token")
	flag.Parse()
	if serverToken == nil || *serverToken == "" {
		log.Fatal("server-token is required")
	}

	pm := postmark.NewClient(*serverToken, "")

	fd, err := ioutil.ReadFile("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	receivers := strings.Split(string(fd), "\n")

	// receivers = []string{
	// 	"gustav@westling.dev",
	// }

	subject := "🐣 This week at Sturdy #18 – What's new in Sturdy v1.8.0"

	for _, receiver := range receivers {
		receiver = strings.TrimSpace(receiver)
		if receiver == "" {
			continue
		}
		log.Println("Sending to", receiver)
		newsletter.Send(pm, subject, "output/2022-05-03.html", receiver)
		time.Sleep(time.Second / 10)
	}
}
