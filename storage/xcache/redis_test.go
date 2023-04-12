package cache

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	testClient = NewRedisClient()
)

func TestMain(m *testing.M) {
	gofakeit.Seed(time.Now().UnixNano())
	m.Run()
	os.Exit(0)
}

func TestDefaultRedisClient_Del(t *testing.T) {
	_, err := testClient.Del("test")
	require.Nil(t, err)
}

func TestDefaultRedisClient_Eval(t *testing.T) {
	r, err := testClient.Eval("return ARGV[1]", nil, []string{"hello"})
	require.Nil(t, err)
	s, ok := r.([]byte)
	require.True(t, ok)
	require.Equal(t, "hello", string(s))
}

func TestDefaultRedisClient_EvalSha(t *testing.T) {
	_, err1 := testClient.Eval("return ARGV[1]", nil, []string{"hello"})
	require.Nil(t, err1)

	r, err2 := testClient.EvalSha("098e0f0d1448c0a81dafe820f66d460eb09263da", nil, []string{"hello"})
	require.Nil(t, err2)
	s, ok := r.([]byte)
	require.True(t, ok)
	require.Equal(t, "hello", string(s))
}

func TestDefaultRedisClient_Get(t *testing.T) {
	randomKey := gofakeit.LetterN(10)
	randomValue := gofakeit.LetterN(10)
	_, err := testClient.Set(randomKey, randomValue)
	require.Nil(t, err)
	v, err := testClient.Get(randomKey)
	require.Nil(t, err)
	require.Equal(t, randomValue, v)
}

func TestDefaultRedisClient_Incr(t *testing.T) {
	randomKey := gofakeit.LetterN(10)
	_, err := testClient.Incr(randomKey, 1)
	require.Nil(t, err)
}

func TestDefaultRedisClient_Set(t *testing.T) {
	randomKey := gofakeit.LetterN(10)
	randomValue := gofakeit.LetterN(10)
	_, err := testClient.Set(randomKey, randomValue)
	require.Nil(t, err)
}

func TestDefaultRedisClient_SetExpired(t *testing.T) {
	randomKey := gofakeit.LetterN(10)
	randomValue := gofakeit.LetterN(10)
	_, err := testClient.SetExpired(randomKey, randomValue, 10)
	require.Nil(t, err)
}

func TestDefaultRedisClient_SetNx(t *testing.T) {
	randomKey := gofakeit.LetterN(10)
	randomValue := gofakeit.LetterN(10)
	_, err := testClient.SetNx(randomKey, randomValue)
	require.Nil(t, err)
}

func TestDefaultRedisClient_SetNxExpired(t *testing.T) {
	randomKey := gofakeit.LetterN(10)
	randomValue := gofakeit.LetterN(10)
	_, err := testClient.SetNxExpired(randomKey, randomValue, 1)
	require.Nil(t, err)
}

func TestNewRedisClient(t *testing.T) {
	client := NewRedisClient(
		SetHost("localhost"),
		SetPort(6379),
		SetPassword(""),
		SetMaxActive(100),
		SetMaxIdle(100),
		SetIdleTimeout(time.Second),
		SetMaxLifetime(time.Minute),
	)
	require.NotNil(t, client)
}
