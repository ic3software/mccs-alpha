package email

import (
	"fmt"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/spf13/viper"
)

type transaction struct{}

var Transaction = &transaction{}

type emailInfo struct {
	InitiatorEmail,
	InitiatorBusinessName,
	ReceiverEmail,
	ReceiverBusinessName string
}

func (tr *transaction) getEmailInfo(t *types.Transaction) *emailInfo {
	var initiatorEmail, initiatorBusinessName, receiverEmail, receiverBusinessName string
	if t.InitiatedBy == t.FromID {
		initiatorBusinessName = t.FromBusinessName
		initiatorEmail = t.FromEmail
		receiverBusinessName = t.ToBusinessName
		receiverEmail = t.ToEmail
	} else {
		initiatorBusinessName = t.ToBusinessName
		initiatorEmail = t.ToEmail
		receiverBusinessName = t.FromBusinessName
		receiverEmail = t.FromEmail
	}
	return &emailInfo{
		initiatorEmail,
		initiatorBusinessName,
		receiverEmail,
		receiverBusinessName,
	}
}

func (tr *transaction) Initiate(
	transactionType string,
	t *types.Transaction,
) error {
	info := tr.getEmailInfo(t)
	url := viper.GetString("url") + "/pending_transactions"

	var body string
	if transactionType == "send" {
		body = info.InitiatorBusinessName + " wants to send " + fmt.Sprintf(
			"%.2f",
			t.Amount,
		) + " Credits to you. <a href=" + url + ">Click here to review this pending transaction</a>."
	} else {
		body = info.InitiatorBusinessName + " wants to receive " + fmt.Sprintf("%.2f", t.Amount) + " Credits from you. <a href=" + url + ">Click here to review this pending transaction</a>."
	}

	d := emailData{
		receiver:      info.ReceiverBusinessName,
		receiverEmail: info.ReceiverEmail,
		subject:       "OCN Transaction Requiring Your Approval",
		text:          body,
		html:          body,
	}
	err := e.send(d)
	if err != nil {
		return err
	}
	return nil
}

func (tr *transaction) Accept(t *types.Transaction) error {
	info := tr.getEmailInfo(t)

	var body string
	if t.InitiatedBy == t.FromID {
		body = info.ReceiverBusinessName + " has accepted the transaction you initiated for -" + fmt.Sprintf(
			"%.2f",
			t.Amount,
		) + " Credits."
	} else {
		body = info.ReceiverBusinessName + " has accepted the transaction you initiated for +" + fmt.Sprintf("%.2f", t.Amount) + " Credits."
	}

	d := emailData{
		receiver:      info.InitiatorBusinessName,
		receiverEmail: info.InitiatorEmail,
		subject:       "OCN Transaction Accepted",
		text:          body,
		html:          body,
	}
	err := e.send(d)
	if err != nil {
		return err
	}
	return nil
}

func (tr *transaction) Cancel(t *types.Transaction, reason string) error {
	info := tr.getEmailInfo(t)

	var body string
	if t.InitiatedBy == t.FromID {
		body = info.InitiatorBusinessName + " has cancelled the transaction it initiated for +" + fmt.Sprintf(
			"%.2f",
			t.Amount,
		) + " Credits."
	} else {
		body = info.InitiatorBusinessName + " has cancelled the transaction it initiated for -" + fmt.Sprintf("%.2f", t.Amount) + " Credits."
	}

	if reason != "" {
		body += "<br/><br/> Reason: <br/><br/>" + reason
	}

	d := emailData{
		receiver:      info.ReceiverBusinessName,
		receiverEmail: info.ReceiverEmail,
		subject:       "OCN Transaction Cancelled",
		text:          body,
		html:          body,
	}
	err := e.send(d)
	if err != nil {
		return err
	}
	return nil
}

func (tr *transaction) CancelBySystem(
	t *types.Transaction,
	reason string,
) error {
	info := tr.getEmailInfo(t)
	body := "The system has cancelled the transaction you initiated with " + info.ReceiverBusinessName + " for the following reason: " + reason
	d := emailData{
		receiver:      info.InitiatorBusinessName,
		receiverEmail: info.InitiatorEmail,
		subject:       "OCN Transaction Cancelled",
		text:          body,
		html:          body,
	}
	err := e.send(d)
	if err != nil {
		return err
	}
	return nil
}

func (tr *transaction) Reject(t *types.Transaction) error {
	info := tr.getEmailInfo(t)

	var body string
	if t.InitiatedBy == t.FromID {
		body = info.ReceiverBusinessName + " has rejected the transaction you initiated for -" + fmt.Sprintf(
			"%.2f",
			t.Amount,
		) + " Credits."
	} else {
		body = info.ReceiverBusinessName + " has rejected the transaction you initiated for +" + fmt.Sprintf("%.2f", t.Amount) + " Credits."
	}

	d := emailData{
		receiver:      info.InitiatorBusinessName,
		receiverEmail: info.InitiatorEmail,
		subject:       "OCN Transaction Rejected",
		text:          body,
		html:          body,
	}
	err := e.send(d)
	if err != nil {
		return err
	}
	return nil
}
