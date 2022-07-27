package util

import "golang.org/x/crypto/bcrypt"

func Hashing(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckHash(text, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err == nil
}
