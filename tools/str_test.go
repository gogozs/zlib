package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStringToArr(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{
			name: "snake",
			s:    "t_user_profile",
			want: []string{"t", "user", "profile"},
		},
		{
			name: "lower camel case",
			s:    "tUserProfile",
			want: []string{"t", "User", "Profile"},
		},
		{
			name: "upper camel case",
			s:    "NewUserProfile",
			want: []string{"New", "User", "Profile"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseWords(tt.s)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToLowerCamelCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"snake", "user", "user"},
		{"snake", "t_user_profile", "tUserProfile"},
		{"upper", "TUserProfile", "tUserProfile"},
		{"lower", "tUserProfile", "tUserProfile"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLowerCamelCase(tt.s)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToUpperCamelCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"snake", "user", "User"},
		{"snake", "t_user_profile", "TUserProfile"},
		{"upper", "TUserProfile", "TUserProfile"},
		{"lower", "tUserProfile", "TUserProfile"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUpperCamelCase(tt.s)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToLowerSnakeString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"snake", "user", "user"},
		{"snake", "t_user_profile", "t_user_profile"},
		{"upper", "TUserProfile", "t_user_profile"},
		{"lower", "tUserProfile", "t_user_profile"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLowerSnakeString(tt.s)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToUpperSnakeString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"snake", "user", "USER"},
		{"snake", "t_user_profile", "T_USER_PROFILE"},
		{"upper", "TUserProfile", "T_USER_PROFILE"},
		{"lower", "tUserProfile", "T_USER_PROFILE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUpperSnakeString(tt.s)
			require.Equal(t, tt.want, got)
		})
	}
}
