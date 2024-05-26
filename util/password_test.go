package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func testPasswords(t *testing.T) {
	passwords := RandomString(6)

	hashPassword1, err := HashPassword(passwords)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword1)

	err = CheckPassword(passwords, hashPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashPassword2, err := HashPassword(passwords)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword2)
	require.NotEmpty(t,hashPassword1,hashPassword2)
}
