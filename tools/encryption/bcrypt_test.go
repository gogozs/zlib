package encryption

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComparePassword(t *testing.T) {
	pwd := "password"
	hashedPwd, err := EncryptPassword("password")
	fmt.Println(hashedPwd)
	require.Nil(t, err)
	require.True(t, ComparePassword(hashedPwd, pwd))
}
