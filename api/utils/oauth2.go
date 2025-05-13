package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
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

	userClaims := UserClaims{
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

func ValidateJWT(token string) (*UserClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	return parsedToken.Claims.(*UserClaims), err
}
