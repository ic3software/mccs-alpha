package bcrypt

import "golang.org/x/crypto/bcrypt"

// Hash hashes the password.
func Hash(password string) (string, error) {
	bytePwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareHash compares the hashedPassword with plainPassword.
func CompareHash(hashedPassword string, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return err
	}
	return nil
}
