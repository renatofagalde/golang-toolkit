package toolkit

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordCrypt interface {
	HashPassword(password string) (string, *RestErr)
	CheckPassword(hashedPassword string, password string) error
}

type Password struct {
}

func NewPassword() PasswordCrypt {
	return &Password{}
}

func (t *Password) HashPassword(password string) (string, *RestErr) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", (&RestErr{}).NewBadRequestError("failed to hash password")
	}

	return string(hashPassword), nil
}

func (t *Password) CheckPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
