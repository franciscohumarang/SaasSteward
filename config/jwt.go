package config

import (
	//"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"saasmanagement/models"
	"time"
)

var jwtSecret = []byte(GetEnv("JWT_SECRET"))

func GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["expires"] = time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyToken(tokenString string) (*models.Claims, error) {
	//sEnc := base64.StdEncoding.EncodeToString(jwtSecret)
	//token, err := jwt.Parse(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
	//dec, _ := base64.StdEncoding.DecodeString(sEnc)
	//	return jwtSecret, nil

	//})

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {

		return nil, err

	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !claims.AccessToken || claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("invalid tokentytyty")
	} else {
		return claims, nil
	}
}
