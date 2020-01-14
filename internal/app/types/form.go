package types

type RegisterData struct {
	User             *User
	Business         *BusinessData
	ConfirmPassword  string
	ConfirmEmail     string
	Terms            string
	RecaptchaSitekey string
}

type UpdateAccountData struct {
	User            *User
	Business        *BusinessData
	Balance         *BalanceLimit
	CurrentPassword string
	ConfirmPassword string
}
