package jwtclaims

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ExpiresAt *jwt.NumericDate `json:"exp"`
	IssuedAt  *jwt.NumericDate `json:"iat"`
	NotBefore *jwt.NumericDate `json:"nbf"`
	Issuer    string           `json:"iss"`
	Subject   string           `json:"sub"`
	Audience  jwt.ClaimStrings `json:"aud"`
	UserID    int64            `json:"uid"`
}

func (cl *Claims) GetExpirationTime() (*jwt.NumericDate, error) { return cl.ExpiresAt, nil }
func (cl *Claims) GetIssuedAt() (*jwt.NumericDate, error)       { return cl.IssuedAt, nil }
func (cl *Claims) GetNotBefore() (*jwt.NumericDate, error)      { return cl.NotBefore, nil }
func (cl *Claims) GetIssuer() (string, error)                   { return cl.Issuer, nil }
func (cl *Claims) GetSubject() (string, error)                  { return cl.Subject, nil }
func (cl *Claims) GetAudience() (jwt.ClaimStrings, error)       { return cl.Audience, nil }
func (cl *Claims) GetUserID() (int64, error)                    { return cl.UserID, nil }
