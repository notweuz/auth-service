package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	jwt.StandardClaims
	Subject         uint64 `json:"sub"`
	Exp             int64  `json:"exp"`
	PasswordVersion uint64 `json:"p_ver"`
}
