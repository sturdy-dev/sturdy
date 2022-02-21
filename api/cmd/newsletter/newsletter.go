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
		"gustav@getsturdy.com",
	}

	subject := "This week at Sturdy - Sturdy is now open-source"

	for _, receiver := range receivers {
		receiver = strings.TrimSpace(receiver)
		log.Println("Sending to", receiver)
		newsletter.Send(sess, subject, "output/2022-02-21.html", receiver)
		time.Sleep(time.Second)
	}
}
