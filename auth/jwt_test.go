package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const jwtSecret = "testSecret"

func TestNewJwtManager(t *testing.T) {
	manager := NewJwtManager(jwtSecret)
	require.NotNil(t, manager)
}

func TestWithTokenDuration(t *testing.T) {
	m := NewJwtManager(jwtSecret, WithTokenDuration(time.Hour))
	require.Equal(t, time.Hour, m.tokenDuration)
}

func TestJwtManager_GenerateToken(t *testing.T) {
	manager := NewJwtManager(jwtSecret)
	token, err := manager.GenerateToken(1)
	require.Nil(t, err)
	require.NotNil(t, token)
}

func TestJwtManager_ParseToken(t *testing.T) {
	var userID uint64 = 1
	manager := NewJwtManager(jwtSecret)
	token, err := manager.GenerateToken(userID)
	require.Nil(t, err)

	parseToken, err := manager.ParseToken(token)
	require.Nil(t, err)
	require.Equal(t, userID, parseToken.UserID)
}

func TestJwtManager_Verify(t *testing.T) {
	m := NewJwtManager(jwtSecret)
	tests := []struct {
		userID           uint64
		generateToken    bool
		expectedErrorNil bool
		expectedUserID   uint64
	}{
		{100, true, true, 100},
		{100, false, false, 0},
		{101, true, true, 101},
		{101, false, false, 0},
	}
	for _, tt := range tests {
		uuid, _ := uuid.NewUUID()
		token := uuid.String()
		if tt.generateToken {
			token, _ = m.GenerateToken(tt.userID)
		}
		uid, err := m.Verify(token)
		require.Equal(t, tt.expectedErrorNil, err == nil)
		require.Equal(t, tt.expectedUserID, uid)
	}
}
