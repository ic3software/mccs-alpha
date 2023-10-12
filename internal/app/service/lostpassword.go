package service

import (
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/spf13/viper"
)

type lostpassword struct{}

var Lostpassword = &lostpassword{}

func (s *lostpassword) Create(l *types.LostPassword) error {
	err := mongo.LostPassword.Create(l)
	if err != nil {
		return e.Wrap(err, "Create failed")
	}
	return nil
}

func (s *lostpassword) FindByToken(token string) (*types.LostPassword, error) {
	lostPassword, err := mongo.LostPassword.FindByToken(token)
	if err != nil {
		return nil, e.Wrap(err, "FindByToken failed")
	}
	return lostPassword, nil
}

func (s *lostpassword) FindByEmail(email string) (*types.LostPassword, error) {
	lostPassword, err := mongo.LostPassword.FindByEmail(email)
	if err != nil {
		return nil, e.Wrap(err, "FindByEmail failed")
	}
	return lostPassword, nil
}

func (s *lostpassword) SetTokenUsed(token string) error {
	err := mongo.LostPassword.SetTokenUsed(token)
	if err != nil {
		return e.Wrap(err, "SetTokenUsed failed")
	}
	return nil
}

func (s *lostpassword) TokenInvalid(l *types.LostPassword) bool {
	if time.Now().
		Sub(l.CreatedAt).
		Seconds() >=
		viper.GetFloat64(
			"reset_password_timeout",
		) ||
		l.TokenUsed == true {
		return true
	}
	return false
}
