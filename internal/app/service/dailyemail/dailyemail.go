package dailyemail

import (
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/ic3network/mccs-alpha/internal/pkg/email"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Run performs the daily email notification.
func Run() {
	users, err := service.User.FindByDailyNotification()
	if err != nil {
		l.Logger.Error("dailyemail failed", zap.Error(err))
		return
	}

	viper.SetDefault("concurrency_num", 1)
	pool := NewPool(viper.GetInt("concurrency_num"))

	for _, user := range users {
		worker := createEmailWorker(user)
		pool.Run(worker)
	}

	pool.Shutdown()
}

func createEmailWorker(u *types.User) func() {
	return func() {
		matchedTags, err := getMatchTags(u)
		if err != nil {
			l.Logger.Error("dailyemail failed", zap.Error(err))
			return
		}
		if len(matchedTags.MatchedOffers) == 0 && len(matchedTags.MatchedWants) == 0 {
			return
		}
		err = email.SendDailyEmailList(u, matchedTags)
		if err != nil {
			l.Logger.Error("dailyemail failed", zap.Error(err))
		}
		err = service.User.UpdateLastNotificationSentDate(u.ID)
		if err != nil {
			l.Logger.Error("dailyemail failed", zap.Error(err))
		}
	}
}

func getMatchTags(user *types.User) (*types.MatchedTags, error) {
	business, err := service.Business.FindByID(user.CompanyID)
	if err != nil {
		return nil, e.Wrap(err, "getMatchTags failed")
	}

	matchedOffers, err := service.Tag.MatchOffers(helper.GetTagNames(business.Offers), user.LastNotificationSentDate)
	if err != nil {
		return nil, e.Wrap(err, "getMatchTags failed")
	}
	matchedWants, err := service.Tag.MatchWants(helper.GetTagNames(business.Wants), user.LastNotificationSentDate)
	if err != nil {
		return nil, e.Wrap(err, "getMatchTags failed")
	}

	return &types.MatchedTags{
		MatchedOffers: matchedOffers,
		MatchedWants:  matchedWants,
	}, nil
}
