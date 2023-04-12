package encryption

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSha1(t *testing.T) {
	require.Equal(t, "098e0f0d1448c0a81dafe820f66d460eb09263da", Sha1("return ARGV[1]"))
}
