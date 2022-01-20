package main

import (
	"getsturdy.com/api/pkg/newsletter"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	conf := aws.NewConfig().WithRegion("eu-north-1")
	sess, err := session.NewSession(conf)
	if err != nil {
		log.Fatal(err)
	}

	// receiversTxt, err := ioutil.ReadFile("output.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// receivers := strings.Split(strings.TrimSpace(string(receiversTxt)), "\n")

	receivers := []string{
		"gustav@westling.dev",
	}

	subject := "This week at Sturdy - Launching the Sturdy App! 🖥"

	for _, receiver := range receivers {
		receiver = strings.TrimSpace(receiver)
		log.Println("Sending to", receiver)
		newsletter.Send(sess, subject, "output/2021-12-07.html", receiver)
		time.Sleep(time.Second)
	}
}
