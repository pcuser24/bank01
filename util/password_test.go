package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	pass := RandomAlphanumericStr(18)

	hashedPass, err := HashPassword(pass)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass)

	err = VerifyPassword(hashedPass, pass)
	require.NoError(t, err)

	wrongPass := RandomAlphanumericStr(18)
	err = VerifyPassword(hashedPass, wrongPass)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPass1, err := HashPassword(pass)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass1)
	require.NotEqual(t, hashedPass, hashedPass1)
}
