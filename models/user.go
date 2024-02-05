package models

import "github.com/golang-jwt/jwt/v4"

// UserAuthIn represents the input to authenticate to a Hagall server.
type UserAuthIn struct {
	Endpoint  string `json:"endpoint"`
	AppKey    string `json:"-"`
	AppSecret string `json:"-"`
}

type UserAuthResponse struct {
	AccessToken string `json:"access_token"`
}

// The claims to generate a Hagall User JWT token.
type HagallUserClaim struct {
	jwt.RegisteredClaims

	AppKey string `json:"app_key"`
}
