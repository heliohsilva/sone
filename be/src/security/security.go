package security

import "golang.org/x/crypto/bcrypt"

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func Compare(password1 string, password2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
}
