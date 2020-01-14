package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sync"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/log"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type adminTagHandler struct {
	once *sync.Once
}

var AdminTagHandler = newAdminTagHandler()

func newAdminTagHandler() *adminTagHandler {
	return &adminTagHandler{
		once: new(sync.Once),
	}
}

func (a *adminTagHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	a.once.Do(func() {
		adminPrivate.Path("/admin-tags").HandlerFunc(a.adminTagPage()).Methods("GET")
		adminPrivate.Path("/admin-tags/search").HandlerFunc(a.searchAdminTags()).Methods("GET")

		public.Path("/api/admin-tags/list/{prefix}").HandlerFunc(a.list()).Methods("GET")
		adminPrivate.Path("/api/admin-tags").HandlerFunc(a.createAdminTag()).Methods("POST")
		adminPrivate.Path("/api/admin-tags/{id}").HandlerFunc(a.renameAdminTag()).Methods("PUT")
		adminPrivate.Path("/api/admin-tags/{id}").HandlerFunc(a.deleteAdminTag()).Methods("DELETE")
	})
}

func (a *adminTagHandler) SaveAdminTags(adminTags []string) error {
	for _, adminTag := range adminTags {
		err := service.AdminTag.Create(adminTag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *adminTagHandler) adminTagPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/admin-tags")
	return func(w http.ResponseWriter, r *http.Request) {
		t.Render(w, r, nil, nil)
	}
}

func (a *adminTagHandler) createAdminTag() func(http.ResponseWriter, *http.Request) {
	type request struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error("CreateAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		if req.Name == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Please enter the admin tag name"))
			return
		}
		req.Name = helper.FormatAdminTag(req.Name)

		_, err = service.AdminTag.FindByName(req.Name)
		if err == nil {
			l.Logger.Info("Admin tag already exists!")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Admin tag already exists!"))
			return
		}

		err = service.AdminTag.Create(req.Name)
		if err != nil {
			l.Logger.Error("CreateAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			if err != nil {
				l.Logger.Error("log.Admin.CreateAdminTag failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(log.Admin.CreateAdminTag(adminUser, req.Name))
			if err != nil {
				l.Logger.Error("log.Admin.CreateAdminTag failed", zap.Error(err))
			}
		}()

		w.WriteHeader(http.StatusCreated)
	}
}

func (a *adminTagHandler) searchAdminTags() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/admin-tags")
	type formData struct {
		Name string
		Page int
	}
	type response struct {
		FormData formData
		Result   *types.FindAdminTagResult
	}
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			l.Logger.Error("SearchAdminTags failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		f := formData{
			Name: r.URL.Query().Get("name"),
			Page: page,
		}
		res := response{FormData: f}

		if f.Name == "" {
			t.Render(w, r, res, []string{"Please enter the admin tag name"})
			return
		}

		findResult, err := service.AdminTag.FindTags(f.Name, int64(f.Page))
		if err != nil {
			l.Logger.Error("SearchAdminTags failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}
		res.Result = findResult

		t.Render(w, r, res, nil)
	}
}

func (a *adminTagHandler) renameAdminTag() func(http.ResponseWriter, *http.Request) {
	type request struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error("RenameAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		if req.Name == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Please enter the tag name"))
			return
		}
		req.Name = helper.FormatAdminTag(req.Name)

		_, err = service.AdminTag.FindByName(req.Name)
		if err == nil {
			l.Logger.Info("RenameAdminTag failed: Admin tag already exists")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Admin tag already exists!"))
			return
		}

		adminTagID, err := primitive.ObjectIDFromHex(req.ID)
		if err != nil {
			l.Logger.Error("RenameAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		adminTag, err := service.AdminTag.FindByID(adminTagID)
		if err != nil {
			l.Logger.Error("RenameAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("AdminTag not found."))
			return
		}
		oldName := adminTag.Name

		go func() {
			err := service.Business.RenameAdminTag(oldName, req.Name)
			if err != nil {
				l.Logger.Error("RenameAdminTag failed", zap.Error(err))
			}
		}()

		adminTag = &types.AdminTag{
			ID:   adminTagID,
			Name: req.Name,
		}
		err = service.AdminTag.Update(adminTag)
		if err != nil {
			l.Logger.Error("RenameAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			if err != nil {
				l.Logger.Error("log.Admin.ModifyAdminTag failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(log.Admin.ModifyAdminTag(adminUser, oldName, req.Name))
			if err != nil {
				l.Logger.Error("log.Admin.ModifyAdminTag failed", zap.Error(err))
			}
		}()

		w.WriteHeader(http.StatusCreated)
	}
}

func (a *adminTagHandler) deleteAdminTag() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		adminTagID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			l.Logger.Error("DeleteAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		adminTag, err := service.AdminTag.FindByID(adminTagID)
		if err != nil {
			l.Logger.Error("DeleteAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = service.AdminTag.DeleteByID(adminTagID)
		if err != nil {
			l.Logger.Error("DeleteAdminTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go func() {
			err := service.Business.DeleteAdminTags(adminTag.Name)
			if err != nil {
				l.Logger.Error("DeleteAdminTags failed", zap.Error(err))
			}
		}()
		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			if err != nil {
				l.Logger.Error("log.Admin.DeleteAdminTag failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(log.Admin.DeleteAdminTag(adminUser, adminTag.Name))
			if err != nil {
				l.Logger.Error("log.Admin.DeleteAdminTag failed", zap.Error(err))
			}
		}()

		w.WriteHeader(http.StatusOK)
	}
}

func (a *adminTagHandler) list() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		prefix := vars["prefix"]

		tags, err := service.AdminTag.TagStartWith(prefix)
		if err != nil {
			l.Logger.Error("controller.AdminTagHandler.List failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(tags)
		if err != nil {
			l.Logger.Error("controller.AdminTagHandler.List failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
