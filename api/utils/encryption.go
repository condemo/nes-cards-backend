package utils

import "golang.org/x/crypto/bcrypt"

func PassEncrypt(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), 14)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func PassVerify(pass, hashedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
	return err == nil
}
