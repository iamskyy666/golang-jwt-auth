package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwt related ops. - create token, parse token, etc.

type Claims struct{
	jwt.RegisteredClaims

	Role string `json:"role"`
}

func CreateToken(jwtSecret string, userID string, role string)(string,error){
	now:=time.Now().UTC()
	expiry:=now.Add(7 * 24 * time.Hour)

	claims:=Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
		Role:role, // For RBA
	}

	token:=jwt.NewWithClaims(jwt.SigningMethodES256,claims)

	signed,err:=token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("⚠️ Token signing failed: %w",err)
	}

	// if all pk, then..
	return signed,nil	
}