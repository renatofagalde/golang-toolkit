package toolkit

import (
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
}

func (t *Password) HashPassword(password string) (string, *RestErr) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//if err != nil {
	//	return "", fmt.Errorf("failed to has password: %w", err)
	//}
	if err != nil {
		return "", (&RestErr{}).NewBadRequestError("failed to hash password")
	}

	return string(hashPassword), nil
}

func (t *Password) CheckPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
