package newsletter

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"strings"

	"github.com/keighl/postmark"
)

func Send(client *postmark.Client, subject string, contentPath string, to string) {
	content, err := ioutil.ReadFile(contentPath)
	if err != nil {
		log.Fatal(err)
	}

	replacer := strings.NewReplacer(
		"BASE64_ENCODED_EMAIL", base64.URLEncoding.EncodeToString([]byte(to)),
	)
	emailContent := replacer.Replace(string(content))

	email := postmark.Email{
		From:          "Gustav at Sturdy <gustav@getsturdy.com>",
		ReplyTo:       "support@getsturdy.com",
		To:            to,
		Subject:       subject,
		HtmlBody:      emailContent,
		Tag:           "newsletter",
		MessageStream: "broadcast",
	}

	res, err := client.SendEmail(email)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("res: %+v", res)
}
