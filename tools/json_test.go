package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToJsonStringWithMaxLength(t *testing.T) {
	text := "abc"
	t.Run("truncation", func(t *testing.T) {
		s := ToJsonStringWithMaxLen(text, 1)
		require.Equal(t, 1, len([]byte(s)))
	})

	t.Run("no truncation", func(t *testing.T) {
		s := ToJsonStringWithMaxLen(text, 6)
		require.Equal(t, 5, len([]byte(s)))
	})
}
