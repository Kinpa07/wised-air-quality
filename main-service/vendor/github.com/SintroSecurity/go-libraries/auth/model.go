package auth

import "github.com/golang-jwt/jwt/v4"

type Token struct {
	jwt.RegisteredClaims
	SkipResourceOwnershipCheck bool
}
