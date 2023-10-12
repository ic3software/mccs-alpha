package jwt

import (
	"crypto/rsa"
	"errors"
	"log"
	"os"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// JWTManager manages JWT operations.
type JWTManager struct {
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

// NewJWTManager initializes and returns a new JWTManager instance.
func NewJWTManager() *JWTManager {
	privateKeyPEM := getEnvOrFallback("jwt.private_key", "JWT_PRIVATE_KEY")
	publicKeyPEM := getEnvOrFallback("jwt.public_key", "JWT_PUBLIC_KEY")

	signKey, err := jwtlib.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		log.Fatal(err)
	}

	verifyKey, err := jwtlib.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		log.Fatal(err)
	}

	return &JWTManager{
		signKey:   signKey,
		verifyKey: verifyKey,
	}
}

type userClaims struct {
	jwtlib.RegisteredClaims
	UserID string `json:"userID"`
	Admin  bool   `json:"admin"`
}

// GenerateToken generates a JWT token for a user.
func (jm *JWTManager) Generate(
	userID string,
	isAdmin bool,
) (string, error) {
	claims := userClaims{
		UserID: userID,
		Admin:  isAdmin,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	return token.SignedString(jm.signKey)
}

// Validate validates a JWT token and returns the associated claims.
func (jm *JWTManager) Validate(tokenString string) (*userClaims, error) {
	claims := &userClaims{}
	token, err := jwtlib.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwtlib.Token) (interface{}, error) {
			return jm.verifyKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func getEnvOrFallback(viperKey, envKey string) string {
	value := viper.GetString(viperKey)
	if value == "" {
		value = os.Getenv(viperKey)
	}
	if value == "" {
		value = os.Getenv(envKey)
	}
	return value
}
