package e

import "github.com/pkg/errors"

type Error struct {
	Code          int
	CustomMessage string
	SystemErr     error
}

func New(code int, systemErr interface{}) error {
	err := Error{
		Code: code,
	}
	if v, ok := systemErr.(error); ok {
		err.SystemErr = v
		return err
	}
	if v, ok := systemErr.(string); ok {
		err.SystemErr = errors.New(v)
		return err
	}
	err.SystemErr = errors.New("Undefined Error")
	return err
}

func CustomMessage(message string) error {
	err := Error{CustomMessage: message}
	err.SystemErr = errors.New("custom error")
	return err
}

func (d Error) Error() string {
	return d.SystemErr.Error()
}

func (d Error) Message() string {
	if d.CustomMessage != "" {
		return d.CustomMessage
	} else if msg, ok := Msg[d.Code]; ok {
		return msg
	}
	return Msg[InternalServerError]
}

func Wrap(err error, message string) error {
	if v, ok := err.(Error); ok {
		return Error{
			Code:      v.Code,
			SystemErr: errors.Wrap(v.SystemErr, message),
		}
	}
	return errors.Wrap(err, message)
}

func IsPasswordInvalid(err error) bool {
	if v, ok := err.(Error); ok {
		return v.Code == PasswordIncorrect
	}
	return false
}

func IsUserNotFound(err error) bool {
	if v, ok := err.(Error); ok {
		return v.Code == UserNotFound
	}
	return false
}
