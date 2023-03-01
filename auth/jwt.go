package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	JwtManager struct {
		jwtSecret     []byte
		tokenDuration time.Duration
	}
	JwtOption struct {
		tokenDuration time.Duration
	}
	JwtOptionFunc func(o *JwtOption)
)

func NewJwtManager(jwtSecret string, opts ...JwtOptionFunc) *JwtManager {
	jwtOption := &JwtOption{
		tokenDuration: time.Hour * 72,
	}
	for _, opt := range opts {
		opt(jwtOption)
	}
	return &JwtManager{
		jwtSecret:     []byte(jwtSecret),
		tokenDuration: jwtOption.tokenDuration,
	}
}

func WithTokenDuration(d time.Duration) JwtOptionFunc {
	return func(o *JwtOption) {
		o.tokenDuration = d
	}
}

func (j *JwtManager) Verify(token string) (uint64, error) {
	chain, err := j.ParseToken(token)
	if err != nil {
		return 0, err
	}
	return chain.UserID, nil
}

type Claims struct {
	UserID uint64
	jwt.StandardClaims
}

func (j *JwtManager) GenerateToken(userID uint64) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.tokenDuration)
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Issuer:    "zlib",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(j.jwtSecret)
	return token, err
}

func (j *JwtManager) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims == nil {
		return nil, errors.Errorf("invalid jwt token")
	}
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, errors.Errorf("invalid jwt token")
}
