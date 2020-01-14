package jwt

import (
	"crypto/rsa"
	"errors"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ic3network/mccs-alpha/global"
	"github.com/spf13/viper"
)

var j *JWT

func init() {
	global.Init()
	j = New()
}

// JWT is a prioritized configuration registry.
type JWT struct {
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

// New returns an initialized JWT instance.
func New() *JWT {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(viper.GetString("jwt.private_key")))
	if err != nil {
		log.Fatal(err)
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(viper.GetString("jwt.public_key")))
	if err != nil {
		log.Fatal(err)
	}

	j := new(JWT)
	j.signKey = signKey
	j.verifyKey = verifyKey
	return j
}

type claims struct {
	jwt.StandardClaims
	UserID string `json:"userID"`
	Admin  bool   `json:"admin"`
}

// GenerateToken generates a jwt token.
func GenerateToken(id string, admin bool) (string, error) { return j.generateToken(id, admin) }
func (j *JWT) generateToken(id string, admin bool) (string, error) {
	c := claims{
		UserID: id,
		Admin:  admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	tokenString, err := token.SignedString(j.signKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken validate a jwt token.
func ValidateToken(tokenString string) (*claims, error) { return j.validateToken(tokenString) }
func (j *JWT) validateToken(tokenString string) (*claims, error) {
	c := &claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (interface{}, error) {
		return j.verifyKey, nil
	})
	if err != nil {
		return &claims{}, err
	}
	if !tkn.Valid {
		return &claims{}, errors.New("Invalid token")
	}
	return c, nil
}
