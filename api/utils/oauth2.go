package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserIDClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"userID"`
}

func CreateJWT(id int64) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	expire, err := strconv.Atoi(os.Getenv("JWT_EXP_DAYS"))
	if err != nil {
		return "", nil
	}
	expireDays := time.Now().Add(time.Hour * 24 * time.Duration(expire))

	userClaims := UserIDClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireDays),
		},
		id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func CreateRefreshJWT(id int64) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	expire, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXP_DAYS"))
	if err != nil {
		return "", nil
	}
	expireDays := time.Now().Add(time.Hour * 24 * time.Duration(expire))

	userClaims := UserIDClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireDays),
		},
		id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(token string) (*UserIDClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &UserIDClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(*UserIDClaims), nil
}
