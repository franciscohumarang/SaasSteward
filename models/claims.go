package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserId       string `json:"user_id"`
	AccessToken  bool   `json:"accessToken"`
	RefreshToken bool   `json:"refreshToken"`
	jwt.StandardClaims
}
