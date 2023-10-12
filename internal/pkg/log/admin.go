package log

import (
	"fmt"
	"strings"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/util"
)

type admin struct{}

var Admin = admin{}

func (a admin) LoginSuccess(
	admin *types.AdminUser,
	ip string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID: admin.ID,
		Email:  admin.Email,
		Action: "admin user login successful",
		// [email] - [IP address]
		ActionDetails: admin.Email + " - " + ip,
		Category:      "admin",
	}
}

func (a admin) LoginFailure(
	admin *types.AdminUser,
	ip string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID: admin.ID,
		Email:  admin.Email,
		Action: "admin user login failed",
		// [email] - [IP address]
		ActionDetails: admin.Email + " - " + ip,
		Category:      "admin",
	}
}

func (a admin) ModifyBusiness(
	admin *types.AdminUser,
	user *types.User,
	oldBusiness *types.Business,
	newBusiness *types.BusinessData,
	oldBalance *types.BalanceLimit,
	newBalance *types.BalanceLimit,
) *types.UserAction {
	modifiedFields := util.CheckDiff(oldBusiness, newBusiness, nil)
	if !helper.SameTags(newBusiness.Offers, oldBusiness.Offers) {
		modifiedFields = append(
			modifiedFields,
			"offers: "+strings.Join(
				helper.GetTagNames(oldBusiness.Offers),
				" ",
			)+" -> "+strings.Join(
				helper.GetTagNames(newBusiness.Offers),
				" ",
			),
		)
	}
	if !helper.SameTags(newBusiness.Wants, oldBusiness.Wants) {
		modifiedFields = append(
			modifiedFields,
			"wants: "+strings.Join(
				helper.GetTagNames(oldBusiness.Wants),
				" ",
			)+" -> "+strings.Join(
				helper.GetTagNames(newBusiness.Wants),
				" ",
			),
		)
	}
	if strings.Join(
		newBusiness.AdminTags,
		" ",
	) != strings.Join(
		oldBusiness.AdminTags,
		" ",
	) {
		modifiedFields = append(
			modifiedFields,
			"adminTags: "+strings.Join(
				oldBusiness.AdminTags,
				" ",
			)+" -> "+strings.Join(
				newBusiness.AdminTags,
				" ",
			),
		)
	}
	modifiedFields = append(
		modifiedFields,
		util.CheckDiff(oldBalance, newBalance, map[string]bool{})...)
	if len(modifiedFields) == 0 {
		return nil
	}
	return &types.UserAction{
		UserID: user.ID,
		Email:  user.Email,
		Action: "admin modified business details",
		ActionDetails: admin.Email + " - " + user.Email + " - " + strings.Join(
			modifiedFields,
			", ",
		),
		Category: "admin",
	}
}

func (a admin) ModifyUser(
	admin *types.AdminUser,
	oldUser *types.User,
	newUser *types.User,
) *types.UserAction {
	modifiedFields := util.CheckDiff(oldUser, newUser, map[string]bool{
		"CurrentLoginIP": true,
		"Password":       true,
		"LastLoginIP":    true,
	})
	if len(modifiedFields) == 0 {
		return nil
	}
	return &types.UserAction{
		UserID: oldUser.ID,
		Email:  newUser.Email,
		Action: "admin modified user details",
		ActionDetails: admin.Email + " - " + newUser.Email + ": " + strings.Join(
			modifiedFields,
			", ",
		),
		Category: "admin",
	}
}

func (a admin) CreateTag(
	admin *types.AdminUser,
	tagName string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID:        admin.ID,
		Email:         admin.Email,
		Action:        "admin created new tag",
		ActionDetails: admin.Email + " - " + tagName,
		Category:      "admin",
	}
}

func (a admin) ModifyTag(
	admin *types.AdminUser,
	old string,
	new string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID:        admin.ID,
		Email:         admin.Email,
		Action:        "admin modified a tag",
		ActionDetails: admin.Email + " - " + old + " -> " + new,
		Category:      "admin",
	}
}

func (a admin) DeleteTag(
	admin *types.AdminUser,
	tagName string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID:        admin.ID,
		Email:         admin.Email,
		Action:        "admin deleted a tag",
		ActionDetails: admin.Email + " - " + tagName,
		Category:      "admin",
	}
}

func (a admin) CreateAdminTag(
	admin *types.AdminUser,
	tagName string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID:        admin.ID,
		Email:         admin.Email,
		Action:        "admin created new admin tag",
		ActionDetails: admin.Email + " - " + tagName,
		Category:      "admin",
	}
}

func (a admin) ModifyAdminTag(
	admin *types.AdminUser,
	old string,
	new string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID:        admin.ID,
		Email:         admin.Email,
		Action:        "admin modified an admin tag",
		ActionDetails: admin.Email + " - " + old + " -> " + new,
		Category:      "admin",
	}
}

func (a admin) DeleteAdminTag(
	admin *types.AdminUser,
	tagName string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID:        admin.ID,
		Email:         admin.Email,
		Action:        "admin deleted an admin tag",
		ActionDetails: admin.Email + " - " + tagName,
		Category:      "admin",
	}
}

func (a admin) Transfer(
	admin *types.AdminUser,
	fromEmail string,
	toEmail string,
	amount float64,
	desc string,
) *types.UserAction {
	admin.Email = strings.ToLower(admin.Email)
	return &types.UserAction{
		UserID: admin.ID,
		Email:  admin.Email,
		Action: "admin transfer for user",
		// admin - [from] -> [to] - [amount]
		ActionDetails: admin.Email + " - " + fromEmail + " -> " + toEmail + " - " + fmt.Sprintf(
			"%.2f",
			amount,
		) + " - " + desc,
		Category: "admin",
	}
}
