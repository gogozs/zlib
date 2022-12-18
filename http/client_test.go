package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpClient_Get(t *testing.T) {
	client := NewHttpClient()
	_, err := client.Get("https://baidu.com")
	require.Nil(t, err)
}
