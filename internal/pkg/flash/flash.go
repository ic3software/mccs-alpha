package flash

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/ic3network/mccs-alpha/global/constant"
)

func Info(w http.ResponseWriter, value string) {
	c := &http.Cookie{
		Name:  constant.Flash.Info,
		Value: encode([]byte(value)),
		Path:  "/",
	}
	http.SetCookie(w, c)
}

func Success(w http.ResponseWriter, value string) {
	c := &http.Cookie{
		Name:  constant.Flash.Success,
		Value: encode([]byte(value)),
		Path:  "/",
	}
	http.SetCookie(w, c)
}

func GetFlash(w http.ResponseWriter, r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return ""
		default:
			return ""
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return ""
	}
	dc := &http.Cookie{
		Name:    name,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
		Path:    "/",
	}
	http.SetCookie(w, dc)
	return string(value)
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
