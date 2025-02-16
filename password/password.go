package password

import (
	"github.com/renatofagalde/golang-toolkit"
	"golang.org/x/crypto/bcrypt"
)

type PasswordCrypt interface {
	HashPassword(password string) (string, *toolkit.RestErr)
	CheckPassword(hashedPassword string, password string) *toolkit.RestErr
}

type Password struct {
}

func NewPassword() PasswordCrypt {
	return &Password{}
}

func (t *Password) HashPassword(password string) (string, *toolkit.RestErr) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", (&toolkit.RestErr{}).NewBadRequestError("failed to hash password")
	}

	return string(hashPassword), nil
}

func (t *Password) CheckPassword(hashedPassword string, password string) *toolkit.RestErr {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		restErr := toolkit.RestErr{}
		causes := []toolkit.Cause{
			{
				Field:   "password",
				Message: err.Error(),
			},
		}
		return restErr.NewRestErr("invalid credentials", "invalid_password", 400, causes)
	}
	return nil
}
