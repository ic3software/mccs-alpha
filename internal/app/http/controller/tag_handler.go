package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/log"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type tagHandler struct {
	once *sync.Once
}

var TagHandler = newTagHandler()

func newTagHandler() *tagHandler {
	return &tagHandler{
		once: new(sync.Once),
	}
}

func (h *tagHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	h.once.Do(func() {
		adminPrivate.Path("/user-tags").HandlerFunc(h.tagPage()).Methods("GET")
		adminPrivate.Path("/user-tags/search").HandlerFunc(h.searchTags()).Methods("GET")

		public.Path("/api/tags/{tagName}").HandlerFunc(h.getTagSuggestions()).Methods("GET")
		adminPrivate.Path("/api/user-tags").HandlerFunc(h.createTag()).Methods("POST")
		adminPrivate.Path("/api/user-tags/{id}").HandlerFunc(h.renameTag()).Methods("PUT")
		adminPrivate.Path("/api/user-tags/{id}").HandlerFunc(h.deleteTag()).Methods("DELETE")
	})
}

func (h *tagHandler) SaveOfferTags(added []string) error {
	for _, tagName := range added {
		// TODO: UpdateOffers
		err := service.Tag.UpdateOffer(tagName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *tagHandler) SaveWantTags(added []string) error {
	for _, tagName := range added {
		// TODO: UpdateWants
		err := service.Tag.UpdateWant(tagName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *tagHandler) tagPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/user-tags")
	return func(w http.ResponseWriter, r *http.Request) {
		t.Render(w, r, nil, nil)
	}
}

func (h *tagHandler) getTagSuggestions() func(http.ResponseWriter, *http.Request) {
	type result struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
		Text  string `json:"text,omitempty"`
	}
	type response struct {
		Success bool     `json:"success,omitempty"`
		Results []result `json:"results,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tagName := vars["tagName"]

		findResult, err := service.Tag.FindTags(tagName, int64(1))
		if err != nil {
			l.Logger.Error("GetTagSuggestions failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		results := make([]result, 0, len(findResult.Tags))
		for _, tag := range findResult.Tags {
			results = append(results, result{
				Name:  tag.Name,
				Value: tag.Name,
				Text:  tag.Name,
			})
		}

		res := response{
			Success: true,
			Results: results,
		}

		js, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (h *tagHandler) createTag() func(http.ResponseWriter, *http.Request) {
	type request struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error("CreateTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		if req.Name == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Please enter the tag name"))
			return
		}

		tagNames := helper.GetTags(req.Name)
		if len(tagNames) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Please enter a valid tag name"))
			return
		}

		tagName := tagNames[0].Name
		_, err = service.Tag.FindByName(tagName)
		if err == nil {
			l.Logger.Info("[CreateTag] failed: Tag already exists")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Tag already exists!"))
			return
		}

		err = service.Tag.Create(tagName)
		if err != nil {
			l.Logger.Error("CreateTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			if err != nil {
				l.Logger.Error("log.Admin.CreateTag failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(log.Admin.CreateTag(adminUser, tagName))
			if err != nil {
				l.Logger.Error("log.Admin.CreateTag failed", zap.Error(err))
			}
		}()

		w.WriteHeader(http.StatusCreated)
	}
}

func (h *tagHandler) searchTags() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/user-tags")
	type formData struct {
		Name string
		Page int
	}
	type response struct {
		FormData formData
		Result   *types.FindTagResult
	}
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			l.Logger.Error("SearchTags failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		f := formData{
			Name: r.URL.Query().Get("name"),
			Page: page,
		}
		res := response{FormData: f}

		if f.Name == "" {
			t.Render(w, r, res, []string{"Please enter the tag name"})
			return
		}

		findResult, err := service.Tag.FindTags(f.Name, int64(f.Page))
		if err != nil {
			l.Logger.Error("SearchTags failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}
		res.Result = findResult

		t.Render(w, r, res, nil)
	}
}

func (h *tagHandler) renameTag() func(http.ResponseWriter, *http.Request) {
	type request struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error("RenameTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		if req.Name == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Please enter the tag name"))
			return
		}

		_, err = service.Tag.FindByName(req.Name)
		if err == nil {
			l.Logger.Info("[RenameTag] failed: Tag already exists")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Tag already exists!"))
			return
		}

		tagID, err := primitive.ObjectIDFromHex(req.ID)
		if err != nil {
			l.Logger.Error("RenameTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		tag, err := service.Tag.FindByID(tagID)
		if err != nil {
			l.Logger.Error("RenameTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Tag not found."))
			return
		}
		oldName := tag.Name

		go func() {
			err := service.Business.RenameTag(oldName, req.Name)
			if err != nil {
				l.Logger.Error("RenameTag failed", zap.Error(err))
			}
		}()

		tag = &types.Tag{
			ID:   tagID,
			Name: req.Name,
		}
		err = service.Tag.Rename(tag)
		if err != nil {
			l.Logger.Error("RenameTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong. Please try again later."))
			return
		}

		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			if err != nil {
				l.Logger.Error("log.Admin.ModifyTag failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(log.Admin.ModifyTag(adminUser, oldName, req.Name))
			if err != nil {
				l.Logger.Error("log.Admin.ModifyTag failed", zap.Error(err))
			}
		}()

		w.WriteHeader(http.StatusCreated)
	}
}

func (h *tagHandler) deleteTag() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		tagID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			l.Logger.Error("DeleteTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tag, err := service.Tag.FindByID(tagID)
		if err != nil {
			l.Logger.Error("DeleteTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = service.Tag.DeleteByID(tagID)
		if err != nil {
			l.Logger.Error("DeleteTag failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go func() {
			err := service.Business.DeleteTag(tag.Name)
			if err != nil {
				l.Logger.Error("DeleteTag failed", zap.Error(err))
			}
		}()
		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			if err != nil {
				l.Logger.Error("log.Admin.DeleteTag failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(log.Admin.DeleteTag(adminUser, tag.Name))
			if err != nil {
				l.Logger.Error("log.Admin.DeleteTag failed", zap.Error(err))
			}
		}()

		w.WriteHeader(http.StatusOK)
	}
}
