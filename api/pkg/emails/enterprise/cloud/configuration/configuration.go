package configuration

type Configuration struct {
	Enable   bool                   `long:"enable" description:"Use AWS SES to send emails"`
	Provider string                 `long:"provider" description:"Which email provider to use" default:"ses"`
	Postmark *PostmarkConfiguration `flags-group:"postmark" namespace:"postmark"`
}

type PostmarkConfiguration struct {
	ServerToken string `long:"server-token" description:"Postmark Server Token"`
}
