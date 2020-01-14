package template

type Data struct {
	User struct {
		ID    string
		Admin bool
	}
	ErrorMessages []string
	Messages      struct {
		Success string
		Info    string
	}
	Yield interface{}
}
