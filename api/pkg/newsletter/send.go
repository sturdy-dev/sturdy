package newsletter

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// PGPASSWORD=$(cat db-pwd) psql -h driva.cqawetpfgboc.eu-north-1.rds.amazonaws.com -U driva -d driva -t -A -F"," -c "select * from users" > output.csv

func Send(session *session.Session, subject string, contentPath string, to string) {
	content, err := ioutil.ReadFile(contentPath)
	if err != nil {
		log.Fatal(err)
	}

	replacer := strings.NewReplacer(
		"BASE64_ENCODED_EMAIL", base64.URLEncoding.EncodeToString([]byte(to)),
	)
	emailContent := replacer.Replace(string(content))

	sesClient := ses.New(session)
	res, err := sesClient.SendEmail(&ses.SendEmailInput{
		Destination:      &ses.Destination{ToAddresses: []*string{&to}},
		ReplyToAddresses: []*string{aws.String("support@getsturdy.com")},
		Source:           aws.String("Gustav at Sturdy <gustav@getsturdy.com>"),
		Message: &ses.Message{
			Subject: &ses.Content{Charset: aws.String("UTF-8"), Data: &subject},
			Body:    &ses.Body{Html: &ses.Content{Charset: aws.String("UTF-8"), Data: &emailContent}},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("res: %+v", res)
}
