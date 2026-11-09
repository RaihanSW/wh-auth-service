package util

import (
	"context"
	"dc-transaction-service/app/model/entity"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(ctx context.Context, param entity.TokenRequest) (entity.AccessToken, error) {
	// Create the JWT claims, which includes the uid and expiry time
	tokenExp := jwt.NewNumericDate(time.Now().Add(time.Duration(param.TokenExp) * time.Hour))
	if param.TokenExp == 0 {
		tokenExp = nil
	}
	claims := entity.Claims{
		UID:      param.UID,
		Username: param.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: tokenExp,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    param.UID,
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(param.Key))
	if err != nil {
		return entity.AccessToken{}, err
	}

	refreshToken := uuid.New()

	return entity.AccessToken{
		Token:              tokenString,
		RefreshToken:       refreshToken.String(),
		RefreshTokenExpire: time.Now().Add(time.Duration(param.RefTokenExp) * time.Hour),
		TokenExpire:        time.Now().Add(time.Duration(param.TokenExp) * time.Hour),
		Claims:             claims,
	}, nil
}

func ValidateToken(ctx context.Context, token string, jwtKey string) (entity.Claims, error) {
	claims := &entity.Claims{}

	// Parse the JWT string and store the result in `claims`.
	// This method will return an error if the token is invalid
	// (expired, or the signature does not match).
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return entity.Claims{}, err
	}
	if !tkn.Valid {
		return entity.Claims{}, errors.New("token is invalid")
	}

	return *claims, nil
}
