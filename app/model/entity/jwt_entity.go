package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UID      string `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AccessToken struct {
	Token              string    `json:"token"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpire time.Time `json:"refresh_token_expire"`
	TokenExpire        time.Time `json:"token_expire"`
	Claims             Claims
}

type TokenRequest struct {
	UID         string
	Username    string
	TokenExp    int // access token expiry in hours
	RefTokenExp int // refresh token expiry in hours
	Key         string
}
