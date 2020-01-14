package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/email"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/ic3network/mccs-alpha/internal/pkg/util"
)

type businessHandler struct {
	once *sync.Once
}

var BusinessHandler = newBusinessHandler()

func newBusinessHandler() *businessHandler {
	return &businessHandler{
		once: new(sync.Once),
	}
}

func (b *businessHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	b.once.Do(func() {
		public.Path("/businesses").HandlerFunc(b.searchBusinessPage()).Methods("GET")
		public.Path("/businesses/search").HandlerFunc(b.searchBusiness()).Methods("GET")
		public.Path("/businessPage/{id}").HandlerFunc(b.businessPage()).Methods("GET")
		private.Path("/businesses/search/match-tags").HandlerFunc(b.searhMatchTags()).Methods("GET")

		private.Path("/api/businessStatus").HandlerFunc(b.businessStatus()).Methods("GET")
		private.Path("/api/getBusinessName").HandlerFunc(b.getBusinessName()).Methods("GET")
		private.Path("/api/tradingMemberStatus").HandlerFunc(b.tradingMemberStatus()).Methods("GET")
		private.Path("/api/contactBusiness").HandlerFunc(b.contactBusiness()).Methods("POST")
	})
}

func (b *businessHandler) FindByID(businessID string) (*types.Business, error) {
	objID, err := primitive.ObjectIDFromHex(businessID)
	if err != nil {
		return nil, err
	}
	business, err := service.Business.FindByID(objID)
	if err != nil {
		return nil, err
	}
	return business, nil
}

func (b *businessHandler) FindByEmail(email string) (*types.Business, error) {
	user, err := service.User.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	bs, err := service.Business.FindByID(user.CompanyID)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (b *businessHandler) FindByUserID(uID string) (*types.Business, error) {
	user, err := UserHandler.FindByID(uID)
	if err != nil {
		return nil, err
	}
	bs, err := service.Business.FindByID(user.CompanyID)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

type searchBusinessFormData struct {
	TagType               string
	Tags                  []*types.TagField
	CreatedOnOrAfter      string
	Category              string
	ShowUserFavoritesOnly bool
	Page                  int
}

type searchBusinessResponse struct {
	IsUserLoggedIn     bool
	FormData           searchBusinessFormData
	Categories         []string
	Result             *types.FindBusinessResult
	FavoriteBusinesses []primitive.ObjectID
}

func (b *businessHandler) searchBusinessPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("businesses")
	return func(w http.ResponseWriter, r *http.Request) {
		adminTags, err := service.AdminTag.GetAll()
		if err != nil {
			l.Logger.Error("SearchBusinessPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		res := searchBusinessResponse{Categories: helper.GetAdminTagNames(adminTags)}
		_, err = UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			res.IsUserLoggedIn = false
		} else {
			res.IsUserLoggedIn = true
		}
		t.Render(w, r, res, nil)
	}
}

func (b *businessHandler) searchBusiness() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("businesses")
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		page, err := strconv.Atoi(q.Get("page"))
		if err != nil {
			l.Logger.Error("SearchBusiness failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		f := searchBusinessFormData{
			TagType:               q.Get("tag_type"),
			Tags:                  helper.ToSearchTags(q.Get("tags")),
			CreatedOnOrAfter:      q.Get("created_on_or_after"),
			Category:              q.Get("category"),
			ShowUserFavoritesOnly: q.Get("show-favorites-only") == "true",
			Page:                  page,
		}
		res := searchBusinessResponse{FormData: f}

		user, err := UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			res.IsUserLoggedIn = false
		} else {
			res.IsUserLoggedIn = true
			res.FavoriteBusinesses = user.FavoriteBusinesses
		}

		c := types.SearchCriteria{
			TagType: f.TagType,
			Tags:    f.Tags,
			Statuses: []string{
				constant.Business.Accepted,
				constant.Trading.Pending,
				constant.Trading.Accepted,
				constant.Trading.Rejected,
			},
			CreatedOnOrAfter:      util.ParseTime(f.CreatedOnOrAfter),
			AdminTag:              f.Category,
			ShowUserFavoritesOnly: f.ShowUserFavoritesOnly,
			FavoriteBusinesses:    res.FavoriteBusinesses,
		}
		findResult, err := service.Business.FindBusiness(&c, int64(f.Page))
		res.Result = findResult
		if err != nil {
			l.Logger.Error("SearchBusiness failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}

		adminTags, err := service.AdminTag.GetAll()
		if err != nil {
			l.Logger.Error("SearchBusiness failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		res.Categories = helper.GetAdminTagNames(adminTags)

		t.Render(w, r, res, nil)
	}
}

func (b *businessHandler) businessPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("business")
	type formData struct {
		IsUserLoggedIn bool
		BusinessEmail  string
		Business       *types.Business
		User           *types.User
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bID := vars["id"]
		business, err := b.FindByID(bID)
		if err != nil {
			l.Logger.Error("BusinessPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		f := formData{
			Business: business,
		}

		businessUser, err := UserHandler.FindByBusinessID(bID)
		if err != nil {
			l.Logger.Error("BusinessPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		f.BusinessEmail = businessUser.Email

		user, err := UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			f.IsUserLoggedIn = false
		} else {
			f.IsUserLoggedIn = true
			f.User = user
		}

		t.Render(w, r, f, nil)
	}
}

func (b *businessHandler) contactBusiness() func(http.ResponseWriter, *http.Request) {
	type request struct {
		BusinessID string `json:"id"`
		Body       string `json:"body"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error("ContactBusiness failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		user, err := UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error("ContactBusiness failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		businessOwner, err := UserHandler.FindByBusinessID(req.BusinessID)
		if err != nil {
			l.Logger.Error("ContactBusiness failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		receiver := businessOwner.FirstName + " " + businessOwner.LastName
		replyToName := user.FirstName + " " + user.LastName
		err = email.SendContactBusiness(receiver, businessOwner.Email, replyToName, user.Email, req.Body)
		if err != nil {
			l.Logger.Error("ContactBusiness failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (b *businessHandler) searhMatchTags() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/businesses/search?"+r.URL.Query().Encode(), http.StatusFound)
	}
}

func (b *businessHandler) businessStatus() func(http.ResponseWriter, *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var business *types.Business
		var err error

		q := r.URL.Query()

		if q.Get("business_id") != "" {
			objID, err := primitive.ObjectIDFromHex(q.Get("business_id"))
			if err != nil {
				l.Logger.Error("BusinessHandler.businessStatus failed", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			business, err = service.Business.FindByID(objID)
			if err != nil {
				l.Logger.Error("BusinessHandler.businessStatus failed", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			business, err = BusinessHandler.FindByUserID(r.Header.Get("userID"))
			if err != nil {
				l.Logger.Error("BusinessHandler.businessStatus failed", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		res := &response{Status: business.Status}
		js, err := json.Marshal(res)
		if err != nil {
			l.Logger.Error("BusinessHandler.businessStatus failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (b *businessHandler) getBusinessName() func(http.ResponseWriter, *http.Request) {
	type response struct {
		Name string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		user, err := service.User.FindByEmail(q.Get("email"))
		if err != nil {
			l.Logger.Error("BusinessHandler.getBusinessName failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		business, err := service.Business.FindByID(user.CompanyID)
		if err != nil {
			l.Logger.Error("BusinessHandler.getBusinessName failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := response{Name: business.BusinessName}
		js, err := json.Marshal(res)
		if err != nil {
			l.Logger.Error("BusinessHandler.getBusinessName failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (b *businessHandler) tradingMemberStatus() func(http.ResponseWriter, *http.Request) {
	type response struct {
		Self  bool `json:"self"`
		Other bool `json:"other"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		objID, err := primitive.ObjectIDFromHex(q.Get("business_id"))
		if err != nil {
			l.Logger.Error("BusinessHandler.tradingMemberStatus failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		other, err := service.Business.FindByID(objID)
		if err != nil {
			l.Logger.Error("BusinessHandler.tradingMemberStatus failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		self, err := BusinessHandler.FindByUserID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error("BusinessHandler.tradingMemberStatus failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := &response{}
		if self.Status == constant.Trading.Accepted {
			res.Self = true
		}
		if other.Status == constant.Trading.Accepted {
			res.Other = true
		}
		js, err := json.Marshal(res)
		if err != nil {
			l.Logger.Error("BusinessHandler.tradingMemberStatus failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
