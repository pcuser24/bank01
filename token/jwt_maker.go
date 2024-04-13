package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

const minSecretKeyLen = 32

// JWTMaker use symmetric key algorithm to sign the token
type JWTMaker struct {
	secreteKey string
}

func NewJWTMaker(secreteKey string) (Maker, error) {
	if len(secreteKey) < minSecretKeyLen {
		return nil, fmt.Errorf("invalid key len: at least %d characters", minSecretKeyLen)
	}
	return &JWTMaker{secreteKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secreteKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidTokenErr
		}
		return []byte(maker.secreteKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ExpiredTokenErr) {
			return nil, ExpiredTokenErr
		}
		return nil, InvalidTokenErr
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, InvalidTokenErr
	}
	return payload, nil
}
