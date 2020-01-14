package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/log"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"github.com/ic3network/mccs-alpha/internal/pkg/util"
	"github.com/ic3network/mccs-alpha/internal/pkg/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/ic3network/mccs-alpha/internal/app/types"
)

type adminBusinessHandler struct {
	once *sync.Once
}

var AdminBusinessHandler = newAdminBusinessHandler()

func newAdminBusinessHandler() *adminBusinessHandler {
	return &adminBusinessHandler{
		once: new(sync.Once),
	}
}

func (a *adminBusinessHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	a.once.Do(func() {
		adminPrivate.Path("/businesses/{id}").HandlerFunc(a.adminBusinessPage()).Methods("GET")
		adminPrivate.Path("/businesses/{id}").HandlerFunc(a.updateBusiness()).Methods("POST")

		adminPrivate.Path("/api/businesses/{id}").HandlerFunc(a.deleteBusiness()).Methods("DELETE")
	})
}

func (a *adminBusinessHandler) adminBusinessPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/business")
	type formData struct {
		Business *types.Business
		Balance  *types.BalanceLimit
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		business, err := BusinessHandler.FindByID(id)
		if err != nil {
			l.Logger.Error("AdminBusinessPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		balance, err := service.BalanceLimit.FindByBusinessID(id)
		if err != nil {
			l.Logger.Error("AdminBusinessPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		f := formData{Business: business, Balance: balance}

		t.Render(w, r, f, nil)
	}
}

func (a *adminBusinessHandler) updateBusiness() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/business")
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		d := helper.GetUpdateData(r)

		vars := mux.Vars(r)
		id := vars["id"]
		bID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			l.Logger.Error("UpdateBusiness failed", zap.Error(err))
			t.Error(w, r, d, err)
			return
		}
		d.Business.ID = bID

		errorMessages := validator.UpdateBusiness(d.Business)
		maxPosBal, err := strconv.ParseFloat(r.FormValue("max_pos_bal"), 64)
		if err != nil {
			errorMessages = append(errorMessages, "Max pos balance should be a number")
		}
		d.Balance.MaxPosBal = math.Abs(maxPosBal)
		maxNegBal, err := strconv.ParseFloat(r.FormValue("max_neg_bal"), 64)
		if err != nil {
			errorMessages = append(errorMessages, "Max neg balance should be a number")
		}
		if math.Abs(maxNegBal) == 0 {
			d.Balance.MaxNegBal = 0
		} else {
			d.Balance.MaxNegBal = math.Abs(maxNegBal)
		}

		// Check if the current balance has exceeded the input balances.
		account, err := service.Account.FindByBusinessID(bID.Hex())
		if err != nil {
			l.Logger.Error("UpdateBusiness failed", zap.Error(err))
			t.Error(w, r, d, err)
			return
		}
		if account.Balance > d.Balance.MaxPosBal {
			errorMessages = append(errorMessages, "The current account balance ("+fmt.Sprintf("%.2f", account.Balance)+") has exceed your max pos balance input")
		}
		if account.Balance < -math.Abs(d.Balance.MaxNegBal) {
			errorMessages = append(errorMessages, "The current account balance ("+fmt.Sprintf("%.2f", account.Balance)+") has exceed your max neg balance input")
		}
		if len(errorMessages) > 0 {
			l.Logger.Error("UpdateBusiness failed", zap.Error(err))
			t.Render(w, r, d, errorMessages)
			return
		}

		// Update Business
		oldBusiness, err := service.Business.FindByID(bID)
		if err != nil {
			l.Logger.Error("UpdateBusiness failed", zap.Error(err))
			t.Error(w, r, d, err)
			return
		}
		offersAdded, offersRemoved := helper.TagDifference(d.Business.Offers, oldBusiness.Offers)
		d.Business.OffersAdded = offersAdded
		d.Business.OffersRemoved = offersRemoved
		wantsAdded, wantsRemoved := helper.TagDifference(d.Business.Wants, oldBusiness.Wants)
		d.Business.WantsAdded = wantsAdded
		d.Business.WantsRemoved = wantsRemoved
		err = service.Business.UpdateBusiness(bID, d.Business, true)
		if err != nil {
			l.Logger.Error("UpdateBusiness failed", zap.Error(err))
			t.Error(w, r, d, err)
			return
		}

		// Update BalanceLimit
		oldBalance, err := service.BalanceLimit.FindByAccountID(account.ID)
		if err != nil {
			l.Logger.Error("UpdateBusiness failed", zap.Error(err))
			t.Error(w, r, d, err)
			return
		}
		err = service.BalanceLimit.Update(account.ID, d.Balance.MaxPosBal, d.Balance.MaxNegBal)
		if err != nil {
			l.Logger.Error("UpdateBusiness failed", zap.Error(err))
			t.Error(w, r, d, err)
			return
		}

		// Update the admin tags collection.
		go func() {
			err := AdminTagHandler.SaveAdminTags(d.Business.AdminTags)
			if err != nil {
				l.Logger.Error("saveAdminTags failed", zap.Error(err))
			}
		}()
		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			user, err := service.User.FindByBusinessID(bID)
			if err != nil {
				l.Logger.Error("log.Admin.ModifyBusiness failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(log.Admin.ModifyBusiness(adminUser, user, oldBusiness, d.Business, oldBalance, d.Balance))
			if err != nil {
				l.Logger.Error("log.Admin.ModifyBusiness failed", zap.Error(err))
			}
		}()

		// Admin Update tags logic:
		// 	1. When a business' status is changed from pending/rejected to accepted.
		// 	   - update all tags.
		// 	2. When the business is in accepted status.
		//	   - only update added tags.
		go func() {
			if !util.IsAcceptedStatus(oldBusiness.Status) && util.IsAcceptedStatus(d.Business.Status) {
				err := service.Business.UpdateAllTagsCreatedAt(oldBusiness.ID, time.Now())
				if err != nil {
					l.Logger.Error("UpdateAllTagsCreatedAt failed", zap.Error(err))
				}
				err = TagHandler.SaveOfferTags(helper.GetTagNames(d.Business.Offers))
				if err != nil {
					l.Logger.Error("saveOfferTags failed", zap.Error(err))
				}
				err = TagHandler.SaveWantTags(helper.GetTagNames(d.Business.Wants))
				if err != nil {
					l.Logger.Error("saveWantTags failed", zap.Error(err))
				}
			}
			if util.IsAcceptedStatus(oldBusiness.Status) && util.IsAcceptedStatus(d.Business.Status) {
				err := TagHandler.SaveOfferTags(d.Business.OffersAdded)
				if err != nil {
					l.Logger.Error("saveOfferTags failed", zap.Error(err))
				}
				err = TagHandler.SaveWantTags(d.Business.WantsAdded)
				if err != nil {
					l.Logger.Error("saveWantTags failed", zap.Error(err))
				}
			}
		}()
		go func() {
			// Set timestamp when first trading status applied.
			if oldBusiness.MemberStartedAt.IsZero() && (oldBusiness.Status == constant.Business.Accepted) && (d.Business.Status == constant.Trading.Accepted) {
				service.Business.SetMemberStartedAt(bID)
			}
		}()

		t.Success(w, r, d, "The business has been updated!")
	}
}

func (a *adminBusinessHandler) deleteBusiness() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		bsID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			l.Logger.Error("DeleteBusiness failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = service.Business.DeleteByID(bsID)
		if err != nil {
			l.Logger.Error("DeleteBusiness failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := service.User.FindByBusinessID(bsID)
		if err != nil {
			l.Logger.Error("DeleteBusiness failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = service.User.DeleteByID(user.ID)
		if err != nil {
			l.Logger.Error("DeleteBusiness failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
