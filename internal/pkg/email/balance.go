package email

import (
	"time"

	"github.com/spf13/viper"
)

type balance struct{}

var Balance = &balance{}

func (b *balance) NonZeroBalance(from time.Time, to time.Time) error {
	body := "Non-zero balance encountered! Please check the timespan from " + from.Format("2006-01-02 15:04:05") + " to " + to.Format("2006-01-02 15:04:05") + " in the posting table."

	d := emailData{
		receiver:      viper.GetString("email_from"),
		receiverEmail: viper.GetString("sendgrid.sender_email"),
		subject:       "[System Check] Non-zero balance encountered",
		text:          body,
		html:          body,
	}
	err := e.send(d)
	if err != nil {
		return err
	}
	return nil
}
