package config

import (
	//"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(GetEnv("JWT_SECRET"))

func GenerateAccessToken(userId string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(userId string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["refresh_token"] = true
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
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

	// Check the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid claims: %v", claims)
		} else {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired: %v", exp)
			} else {
				return claims, nil
			}
		}
	} else {
		return nil, fmt.Errorf("invalid token")
	}

}

func VerifyRefreshToken(tokenString string) (jwt.MapClaims, error) {
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

	// Check the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		isRefreshToken, ok := claims["refresh_token"].(bool)
		if !ok && !isRefreshToken {
			return nil, fmt.Errorf("invalid claims: %v", claims)
		}

		exp, ok := claims["exp"].(float64)
		if !ok {

			return nil, fmt.Errorf("invalid claims: %v", claims)

		} else {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("refresh token expired")
			} else {
				return claims, nil
			}
		}
	} else {
		return nil, fmt.Errorf("invalid token")
	}

}
