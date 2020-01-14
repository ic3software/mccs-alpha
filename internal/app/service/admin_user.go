package service

import (
	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/bcrypt"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type adminUser struct{}

var AdminUser = &adminUser{}

func (a *adminUser) Login(email string, password string) (*types.AdminUser, error) {
	user, err := mongo.AdminUser.FindByEmail(email)
	if err != nil {
		return &types.AdminUser{}, e.Wrap(err, "login admin user failed")
	}

	err = bcrypt.CompareHash(user.Password, password)
	if err != nil {
		return &types.AdminUser{}, e.New(e.PasswordIncorrect, err)
	}

	return user, nil
}

func (a *adminUser) FindByID(id primitive.ObjectID) (*types.AdminUser, error) {
	adminUser, err := mongo.AdminUser.FindByID(id)
	if err != nil {
		return nil, e.Wrap(err, "service.AdminUser.FindByID failed")
	}
	return adminUser, nil
}

func (a *adminUser) FindByEmail(email string) (*types.AdminUser, error) {
	adminUser, err := mongo.AdminUser.FindByEmail(email)
	if err != nil {
		return nil, e.Wrap(err, "service.AdminUser.FindByEmail failed")
	}
	return adminUser, nil
}

func (a *adminUser) UpdateLoginInfo(id primitive.ObjectID, ip string) error {
	loginInfo, err := mongo.AdminUser.GetLoginInfo(id)
	if err != nil {
		return e.Wrap(err, "service.AdminUser.UpdateLoginInfo failed")
	}

	newLoginInfo := &types.LoginInfo{
		CurrentLoginIP: ip,
		LastLoginIP:    loginInfo.CurrentLoginIP,
		LastLoginDate:  loginInfo.CurrentLoginDate,
	}

	err = mongo.AdminUser.UpdateLoginInfo(id, newLoginInfo)
	if err != nil {
		return e.Wrap(err, "AdminUserService UpdateLoginInfo failed")
	}
	return nil
}
