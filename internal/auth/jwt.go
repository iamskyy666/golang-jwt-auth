package auth

import (
	"errors"
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

	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	signed,err:=token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("⚠️ Token signing failed: %w",err)
	}

	// if all ok, then..
	return signed,nil	
}

// for RBA..
func ParseToken(jwtSecret string, tokenStr string)(Claims, error){
	var claims Claims

	parsed,err:=jwt.ParseWithClaims(tokenStr, &claims,
	func(t *jwt.Token) (any, error) {
		// verify the algo.
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg(){
			return nil,fmt.Errorf(" ⚠️ Unexpected signing-method: %v",t.Header["alg"])
		}
		return []byte(jwtSecret),nil
	},
	jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),

	)

	// Checks
	if err!=nil{
		return Claims{},fmt.Errorf(" ⚠️ Token-parsing failed: %w",err)
	}

	if !parsed.Valid{
		return Claims{},errors.New(" ⚠️ Invalid token!")
	}

	if claims.Subject==""{
		return Claims{},errors.New(" ⚠️ Token missing subject!")
	}
	
	// if all ok, then.. ✅
	return claims,nil
}