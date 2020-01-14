package e

const (
	UserNotFound = iota
	BusinessNotFound
	InternalServerError
	EmailExisted
	PasswordIncorrect
	AccountLocked
	TokenInvalid
	InvalidPageNumber
	ExceedMaxPosBalance
	ExceedMaxNegBalance
)

var Msg = map[int]string{
	UserNotFound:        "Email address not found.",
	BusinessNotFound:    "Business not found.",
	EmailExisted:        "Email address is already registered.",
	TokenInvalid:        "Invalid token.",
	PasswordIncorrect:   "Invalid password.",
	AccountLocked:       "Your account has been temporarily locked for 15 minutes. Please try again later.",
	InternalServerError: "Sorry, something went wrong. Please try again later.",
	InvalidPageNumber:   "Invalid page number: should start with 1.",
	ExceedMaxPosBalance: "Transfer rejected: receiver will exceed maximum balance limit.",
	ExceedMaxNegBalance: "Transfer rejected: you will exceed your maximum negative balance limit.",
}
