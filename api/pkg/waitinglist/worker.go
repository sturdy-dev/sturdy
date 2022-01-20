package waitinglist

import (
	"encoding/base64"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/api/gmail/v1"
)

func Worker(logger *zap.Logger, repo WaitingListRepo, gmailService *gmail.Service) {
	for {
		// time.Sleep(time.Second * 10)

		toInvite, err := repo.ToSendInvitesTo()
		if err != nil {
			logger.Error("failed to get users to invite", zap.Error(err))
			continue
		}

		for _, u := range toInvite {
			err := send(logger, gmailService, u.Email)
			if err != nil {
				logger.Error("failed to send invite", zap.Error(err))
				continue
			}

			// Mark as invited
			err = repo.MarkEmailAsInvited(u.Email)
			if err != nil {
				logger.Error("failed to mark user as invited", zap.Error(err))
				continue
			}
		}

		time.Sleep(time.Second * 4)
	}
}

func send(logger *zap.Logger, srv *gmail.Service, to string) error {
	var message gmail.Message

	user := "me"

	msg := []byte(
		fmt.Sprintf("To: %s\r\n", to) +
			fmt.Sprintf("From: %s\r\n", "Gustav Westling <gustav@getsturdy.com>") +
			fmt.Sprintf("Subject: %s\n", "Welcome to Sturdy!") +
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n" +
			fmt.Sprintf(`
Congratulations, you've been invited to the early-access of <a href="https://getsturdy.com">Sturdy</a>!<br><br>

<a href="https://getsturdy.com/auth/%s">Create your account</a> and get started with Sturdy today (it's free).<br><br>

At Sturdy we're building the future of how developers will collaborate on code, without having to deal with git.
Our journey is just getting started, but we're excited to have <b>you</b> as one of our first users.<br><br>

You can read more about Sturdy and how to get started in <a href="http://getsturdy.com/docs">the documentation</a> and from our latest <a href="https://getsturdy.com/blog">blog posts</a>.<br><br>

I'm happy to have you onboard, and if you have any questions, reply to this email and I'll help you out.<br><br>

Have a nice day,<br>
Gustav Westling
`, to))

	message.Raw = base64.URLEncoding.EncodeToString(msg)
	res, err := srv.Users.Messages.Send(user, &message).Do()
	if err != nil {
		return err
	}
	logger.Info("sent invite email", zap.Any("res", res))
	return nil
}
