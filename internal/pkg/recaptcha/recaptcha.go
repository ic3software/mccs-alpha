package recaptcha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ic3network/mccs-alpha/global"
	"github.com/spf13/viper"
)

var r *Recaptcha

func init() {
	global.Init()
	r = New()
}

// Recaptcha is a prioritized configuration registry.
type Recaptcha struct {
	Secret string
	errMsg string
}

// New returns an initialized recaptcha instance.
func New() *Recaptcha {
	j := new(Recaptcha)
	j.Secret = viper.GetString("recaptcha.secret_key")
	return j
}

// Struct for parsing json in google's response
type googleResponse struct {
	Success    bool
	ErrorCodes []string `json:"error-codes"`
}

// url to post submitted re-captcha response to
var postURL = "https://www.google.com/recaptcha/api/siteverify"

// Verify method, verifies if current request have valid re-captcha response and returns true or false
// This method also records any errors in validation.
// These errors can be received by calling LastError() method.
func Verify(req http.Request) bool { return r.verify(req) }
func (r *Recaptcha) verify(req http.Request) bool {
	response := req.FormValue("g-recaptcha-response")
	return r.verifyResponse(response)
}

// VerifyResponse is a method similar to `Verify`; but doesn't parse the form for you. Useful if
// you're receiving the data as a JSON object from a javascript app or similar.
func (r *Recaptcha) verifyResponse(response string) bool {
	if response == "" {
		r.errMsg = "Please select captcha first."
		return false
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.PostForm(
		postURL,
		url.Values{"secret": {r.Secret}, "response": {response}},
	)
	if err != nil {
		r.errMsg = err.Error()
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.errMsg = err.Error()
		return false
	}
	gr := new(googleResponse)
	err = json.Unmarshal(body, gr)
	if err != nil {
		r.errMsg = err.Error()
		return false
	}
	if !gr.Success {
		r.errMsg = gr.ErrorCodes[len(gr.ErrorCodes)-1]
	}
	return gr.Success
}

// Error returns errors occurred in last re-captcha validation attempt
func Error() []string { return r.error() }
func (r *Recaptcha) error() []string {
	return []string{r.errMsg}
}
