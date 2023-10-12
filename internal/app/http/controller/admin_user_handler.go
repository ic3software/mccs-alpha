package controller

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/cookie"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/ip"
	"github.com/ic3network/mccs-alpha/internal/pkg/jwt"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/log"
	"github.com/ic3network/mccs-alpha/internal/pkg/recaptcha"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"github.com/ic3network/mccs-alpha/internal/pkg/validator"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type adminUserHandler struct {
	once *sync.Once
}

var AdminUserHandler = newAdminUserHandler()

func newAdminUserHandler() *adminUserHandler {
	return &adminUserHandler{
		once: new(sync.Once),
	}
}

func (a *adminUserHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	a.once.Do(func() {
		adminPrivate.Path("").HandlerFunc(a.dashboardPage()).Methods("GET")
		adminPublic.Path("/login").HandlerFunc(a.loginPage()).Methods("GET")
		adminPublic.Path("/login").HandlerFunc(a.loginHandler()).Methods("POST")
		adminPrivate.Path("/logout").HandlerFunc(a.logoutHandler()).Methods("GET")
		adminPrivate.Path("/users/{id}").HandlerFunc(a.userPage()).Methods("GET")
		adminPrivate.Path("/users/{id}").HandlerFunc(a.updateUser()).Methods("POST")
	})
}

func (a *adminUserHandler) dashboardPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/dashboard")
	return func(w http.ResponseWriter, r *http.Request) {
		objID, err := primitive.ObjectIDFromHex(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error("AdminDashboardPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		adminUser, err := service.AdminUser.FindByID(objID)
		if err != nil {
			l.Logger.Error("AdminDashboardPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		t.Render(w, r, adminUser, nil)
	}
}

func (a *adminUserHandler) loginPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("/admin/login")
	type formData struct {
		Email            string
		Password         string
		RecaptchaSitekey string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		t.Render(w, r, formData{RecaptchaSitekey: viper.GetString("recaptcha.site_key")}, nil)
	}
}

func (a *adminUserHandler) loginHandler() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("/admin/login")
	type formData struct {
		Email            string
		Password         string
		RecaptchaSitekey string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		f := formData{
			Email:            r.FormValue("email"),
			Password:         r.FormValue("password"),
			RecaptchaSitekey: viper.GetString("recaptcha.site_key"),
		}

		if viper.GetString("env") == "production" {
			isValid := recaptcha.Verify(*r)
			if !isValid {
				l.Logger.Error("AdminLoginHandler failed", zap.Strings("errs", recaptcha.Error()))
				t.Render(w, r, f, recaptcha.Error())
				return
			}
		}

		user, err := service.AdminUser.Login(f.Email, f.Password)
		if err != nil {
			l.Logger.Info("AdminLoginHandler failed", zap.Error(err))
			t.Error(w, r, f, err)
			go func() {
				user, err := service.AdminUser.FindByEmail(f.Email)
				if err != nil {
					if !e.IsUserNotFound(err) {
						l.Logger.Error("BuildLoginFailureAction failed", zap.Error(err))
					}
					return
				}
				err = service.UserAction.Log(log.Admin.LoginFailure(user, ip.FromRequest(r)))
				if err != nil {
					l.Logger.Error("BuildLoginFailureAction failed", zap.Error(err))
				}
			}()
			return
		}

		token, err := jwt.NewJWTManager().GenerateToken(user.ID.Hex(), true)
		http.SetCookie(w, cookie.CreateCookie(token))

		go func() {
			err := service.AdminUser.UpdateLoginInfo(user.ID, ip.FromRequest(r))
			if err != nil {
				l.Logger.Error("AdminLoginHandler failed", zap.Error(err))
			}
		}()
		go func() {
			err := service.UserAction.Log(log.Admin.LoginSuccess(user, ip.FromRequest(r)))
			if err != nil {
				l.Logger.Error("log.Admin.LoginSuccess failed", zap.Error(err))
			}
		}()

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (a *adminUserHandler) logoutHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, cookie.ResetCookie())
		http.Redirect(w, r, "/admin/login", http.StatusFound)
	}
}

func (a *adminUserHandler) userPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/user")
	type formData struct {
		User *types.User
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]
		user, err := UserHandler.FindByID(userID)
		if err != nil {
			l.Logger.Error("UserPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		f := formData{User: user}

		t.Render(w, r, f, nil)
	}
}

func (a *adminUserHandler) updateUser() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/user")
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		userID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			l.Logger.Error("UpdateUser failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		r.ParseForm()
		updateData := helper.GetUpdateData(r)
		updateData.User.ID = userID

		errorMessages := validator.ValidateUser(updateData.User)

		if (r.FormValue("origin_email") != updateData.User.Email) && service.User.UserEmailExists(updateData.User.Email) {
			errorMessages = append(errorMessages, "Email address is already registered")
		}

		if len(errorMessages) > 0 {
			l.Logger.Error("UpdateUser failed", zap.Error(err))
			t.Render(w, r, updateData, errorMessages)
			return
		}

		oldUser, err := service.User.FindByEmail(r.FormValue("origin_email"))
		if err != nil {
			l.Logger.Error("UpdateUser failed", zap.Error(err))
			t.Error(w, r, updateData, err)
			return
		}

		err = service.User.AdminUpdateUser(updateData.User)
		if err != nil {
			l.Logger.Error("UpdateUser failed", zap.Error(err))
			t.Error(w, r, updateData, err)
			return
		}

		if updateData.User.Password != "" || updateData.ConfirmPassword != "" {
			errorMessages := validator.ValidatePassword(updateData.User.Password, updateData.ConfirmPassword)
			if len(errorMessages) > 0 {
				l.Logger.Error("UpdateUser failed", zap.Strings("input invalid", errorMessages))
				t.Render(w, r, updateData, errorMessages)
				return
			}
			err = service.User.ResetPassword(updateData.User.Email, updateData.ConfirmPassword)
			if err != nil {
				l.Logger.Error("UpdateUser failed", zap.Error(err))
				t.Error(w, r, updateData, err)
				return
			}
		}

		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			if err != nil {
				l.Logger.Error("log.Admin.ModifyUser failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(log.Admin.ModifyUser(adminUser, oldUser, updateData.User))
			if err != nil {
				l.Logger.Error("log.Admin.ModifyUser failed", zap.Error(err))
			}
		}()

		t.Success(w, r, updateData, "The user has been updated!")
	}
}
