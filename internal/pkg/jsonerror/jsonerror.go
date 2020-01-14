package jsonerror

type JE struct {
	Code    string
	message string
}

// New creates a new JE struct.
func New(code, message string) JE {
	j := JE{Code: code, message: message}
	return j
}

func (j JE) Render() map[string]string {
	return map[string]string{"code": j.Code, "message": j.message}
}
