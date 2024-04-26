package commons

import (
	"errors"
	"go_clean_architecture/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtProvider struct {
	secret string
}

type TokenPayload struct {
	Id uint
}

type Token struct {
	Token   string
	Expiry  int64
	Created time.Time
}

func NewJwtProvider() *jwtProvider {
	return &jwtProvider{
		secret: config.Env.AUTH_TOKEN,
	}
}

type claimsPayload struct {
	Payload TokenPayload
	jwt.RegisteredClaims
}

func (j jwtProvider) Generate(data TokenPayload, timeDuration time.Duration) (*Token, error) {
	expiry := time.Now().Add(timeDuration).Unix()
	claims := claimsPayload{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiry, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}
	return &Token{
		Token:   tokenStr,
		Expiry:  expiry,
		Created: time.Now(),
	}, nil
}

func (j jwtProvider) Validate(tokenString string) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claimsPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, ErrInternal(err)
	}
	claims, ok := token.Claims.(*claimsPayload)
	if !ok {
		return nil, ErrInternal(errors.New("cannot validate claims"))
	}

	return &claims.Payload, nil
}
