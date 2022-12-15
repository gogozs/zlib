package http

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHttpClient_Get(t *testing.T) {
	client := NewHttpClient()
	_, err := client.Get("https://baidu.com")
	require.Nil(t, err)
}
