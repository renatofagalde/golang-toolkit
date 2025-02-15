package password

import (
	"github.com/renatofagalde/golang-toolkit"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func Test_password(t *testing.T) {
	password := toolkit.RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(hashedPassword, password)
	require.NoError(t, err)

	wrongPassword := toolkit.RandomString(6)

	err = CheckPassword(hashedPassword, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
