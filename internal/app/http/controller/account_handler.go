package controller

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/log"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"github.com/ic3network/mccs-alpha/internal/pkg/util"
	"github.com/ic3network/mccs-alpha/internal/pkg/validator"
	"go.uber.org/zap"

	"github.com/ic3network/mccs-alpha/internal/app/types"
)

type accountHandler struct {
	once *sync.Once
}

var AccountHandler = newAccountHandler()

func newAccountHandler() *accountHandler {
	return &accountHandler{
		once: new(sync.Once),
	}
}

func (a *accountHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	a.once.Do(func() {
		private.Path("/account").HandlerFunc(a.accountPage()).Methods("GET")
		private.Path("/account").HandlerFunc(a.updateAccount()).Methods("POST")
		adminPrivate.Path("/accounts").
			HandlerFunc(a.searchAccountPage()).
			Methods("GET")
		adminPrivate.Path("/accounts/search").
			HandlerFunc(a.searchAccount()).
			Methods("GET")
	})
}

type searchAccountFormData struct {
	TagType          string
	Tags             []*types.TagField
	CreatedOnOrAfter string
	Status           string
	BusinessName     string
	LocationCity     string
	LocationCountry  string
	Category         string
	LastName         string
	Email            string
	Filter           string
	Page             int
}

type account struct {
	Business *types.Business
	User     *types.User
	Balance  float64
}

type findAccountResult struct {
	Accounts        []account
	NumberOfResults int
	TotalPages      int
}

type sreachResponse struct {
	FormData  searchAccountFormData
	AdminTags []*types.AdminTag
	Result    *findAccountResult
}

func (a *accountHandler) searchAccountPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/accounts")
	return func(w http.ResponseWriter, r *http.Request) {
		adminTags, err := service.AdminTag.GetAll()
		if err != nil {
			l.Logger.Error("SearchAccountPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		res := sreachResponse{AdminTags: adminTags}
		t.Render(w, r, res, nil)
	}
}

func (a *accountHandler) searchAccount() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/accounts")
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		page, err := strconv.Atoi(q.Get("page"))
		if err != nil {
			l.Logger.Error("SearchAccount failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		f := searchAccountFormData{
			TagType:          q.Get("tag_type"),
			Tags:             helper.ToSearchTags(q.Get("tags")),
			Status:           q.Get("status"),
			BusinessName:     q.Get("business_name"),
			CreatedOnOrAfter: q.Get("created_on_or_after"),
			LocationCity:     q.Get("location_city"),
			LocationCountry:  q.Get("location_country"),
			Category:         q.Get("category"),
			LastName:         q.Get("last_name"),
			Email:            q.Get("email"),
			Filter:           q.Get("filter"),
			Page:             page,
		}
		res := sreachResponse{FormData: f, Result: new(findAccountResult)}

		if f.Filter != "business" && f.LastName == "" && f.Email == "" {
			t.Render(
				w,
				r,
				res,
				[]string{"Please enter at least one search criteria."},
			)
			return
		}

		adminTags, err := service.AdminTag.GetAll()
		if err != nil {
			l.Logger.Error("SearchAccount failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		res.AdminTags = adminTags

		// Search All Status
		var status []string
		if f.Status == constant.ALL {
			status = []string{
				constant.Business.Pending,
				constant.Business.Accepted,
				constant.Business.Rejected,
				constant.Trading.Pending,
				constant.Trading.Accepted,
				constant.Trading.Rejected,
			}
		} else {
			status = []string{f.Status}
		}

		findResult := new(types.FindBusinessResult)
		if f.Filter == "business" {
			c := types.SearchCriteria{
				TagType:          f.TagType,
				Tags:             f.Tags,
				Statuses:         status,
				BusinessName:     f.BusinessName,
				CreatedOnOrAfter: util.ParseTime(f.CreatedOnOrAfter),
				LocationCity:     f.LocationCity,
				LocationCountry:  f.LocationCountry,
				AdminTag:         f.Category,
			}
			findResult, err = service.Business.FindBusiness(&c, int64(f.Page))
			if err != nil {
				l.Logger.Error("SearchAccount failed", zap.Error(err))
				t.Error(w, r, res, err)
				return
			}
			res.Result.TotalPages = findResult.TotalPages
			res.Result.NumberOfResults = findResult.NumberOfResults
		}

		accounts := make([]account, 0)
		// Find the user and account balance using business id.
		for _, business := range findResult.Businesses {
			user, err := service.User.FindByBusinessID(business.ID)
			if err != nil {
				l.Logger.Error("SearchAccount failed", zap.Error(err))
				t.Error(w, r, res, err)
				return
			}
			acc, err := service.Account.FindByBusinessID(business.ID.Hex())
			if err != nil {
				l.Logger.Error("SearchAccount failed", zap.Error(err))
				t.Error(w, r, res, err)
				return
			}
			accounts = append(accounts, account{
				Business: business,
				User:     user,
				Balance:  acc.Balance,
			})
		}
		res.Result.Accounts = accounts

		if len(res.Result.Accounts) > 0 || f.Filter == "business" {
			t.Render(w, r, res, nil)
			return
		}

		// The logic for searching by user last name and email.
		u := types.User{
			LastName: f.LastName,
			Email:    f.Email,
		}
		findUserResult, err := service.User.FindUsers(&u, int64(f.Page))
		if err != nil {
			l.Logger.Error("SearchAccount failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}
		res.Result.TotalPages = findUserResult.TotalPages
		res.Result.NumberOfResults = findUserResult.NumberOfResults

		// Find the business and account balance.
		for _, user := range findUserResult.Users {
			business, err := service.Business.FindByID(user.CompanyID)
			if err != nil {
				l.Logger.Error("SearchAccount failed", zap.Error(err))
				t.Error(w, r, res, err)
				return
			}
			acc, err := service.Account.FindByBusinessID(business.ID.Hex())
			if err != nil {
				l.Logger.Error("SearchAccount failed", zap.Error(err))
				t.Error(w, r, res, err)
				return
			}
			accounts = append(accounts, account{
				Business: business,
				User:     user,
				Balance:  acc.Balance,
			})
		}
		res.Result.Accounts = accounts

		t.Render(w, r, res, nil)
	}
}

func (a *accountHandler) accountPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("account")
	type request struct {
		User     *types.User
		Business *types.Business
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error("AccountPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		business, err := service.Business.FindByID(user.CompanyID)
		if err != nil {
			l.Logger.Error("AccountPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		t.Render(w, r, request{User: user, Business: business}, nil)
	}
}

func (a *accountHandler) updateAccount() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("account")
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		formData := helper.GetUpdateData(r)

		// Find the user and he's business.
		user, err := service.User.FindByEmail(formData.User.Email)
		if err != nil {
			l.Logger.Error("appServer UpdateAccount failed", zap.Error(err))
			t.Error(w, r, formData, err)
			return
		}
		oldBusiness, err := service.Business.FindByID(user.CompanyID)
		if err != nil {
			l.Logger.Error("appServer UpdateAccount failed", zap.Error(err))
			return
		}

		// Validate the user inputs.
		errorMessages := []string{}
		if formData.CurrentPassword != "" {
			_, err := service.User.Login(
				formData.User.Email,
				formData.CurrentPassword,
			)
			if err != nil {
				l.Logger.Error("appServer UpdateAccount failed", zap.Error(err))
				t.Error(w, r, formData, err)
				return
			}
		}
		errorMessages = validator.Account(formData)
		if oldBusiness.Status == constant.Trading.Accepted {
			// Additional validation if the business status is "tradingAccepted".
			data := helper.Trading.GetUpdateData(r)
			errorMessages = append(errorMessages, data.Validate()...)
		}
		if len(errorMessages) > 0 {
			l.Logger.Info(
				"appServer UpdateAccount failed",
				zap.Strings("input invalid", errorMessages),
			)
			t.Render(w, r, formData, errorMessages)
			return
		}

		formData.User.ID = user.ID
		err = service.User.UpdateUserInfo(formData.User)
		if err != nil {
			l.Logger.Error("appServer UpdateAccount failed", zap.Error(err))
			t.Error(w, r, formData, err)
			return
		}

		offersAdded, offersRemoved := helper.TagDifference(
			formData.Business.Offers,
			oldBusiness.Offers,
		)
		formData.Business.OffersAdded = offersAdded
		formData.Business.OffersRemoved = offersRemoved
		wantsAdded, wantsRemoved := helper.TagDifference(
			formData.Business.Wants,
			oldBusiness.Wants,
		)
		formData.Business.WantsAdded = wantsAdded
		formData.Business.WantsRemoved = wantsRemoved

		err = service.Business.UpdateBusiness(
			user.CompanyID,
			formData.Business,
			false,
		)
		if err != nil {
			l.Logger.Error("appServer UpdateAccount failed", zap.Error(err))
			t.Error(w, r, formData, err)
			return
		}

		if formData.CurrentPassword != "" && formData.ConfirmPassword != "" {
			err = service.User.ResetPassword(
				user.Email,
				formData.ConfirmPassword,
			)
			if err != nil {
				l.Logger.Error("appServer UpdateAccount failed", zap.Error(err))
				t.Error(w, r, formData, err)
				return
			}
		}

		go func() {
			err := service.UserAction.Log(
				log.User.ModifyAccount(
					user,
					formData.User,
					oldBusiness,
					formData.Business,
				),
			)
			if err != nil {
				l.Logger.Error(
					"BuildModifyAccountAction failed",
					zap.Error(err),
				)
			}
		}()

		// User Update tags logic:
		// 	1. Update the tags collection only when the business is in accepted status.
		go func() {
			if util.IsAcceptedStatus(oldBusiness.Status) {
				err := TagHandler.SaveOfferTags(formData.Business.OffersAdded)
				if err != nil {
					l.Logger.Error("saveOfferTags failed", zap.Error(err))
				}
				err = TagHandler.SaveWantTags(formData.Business.WantsAdded)
				if err != nil {
					l.Logger.Error("saveWantTags failed", zap.Error(err))
				}
			}
		}()

		t.Success(w, r, formData, "Your account has been updated!")
	}
}

func (a *accountHandler) FindByUserID(uID string) (*types.Account, error) {
	business, err := BusinessHandler.FindByUserID(uID)
	if err != nil {
		return nil, e.Wrap(err, "controller.Business.FindByUserID failed")
	}
	account, err := service.Account.FindByBusinessID(business.ID.Hex())
	if err != nil {
		return nil, e.Wrap(err, "controller.Business.FindByUserID failed")
	}
	return account, nil
}
