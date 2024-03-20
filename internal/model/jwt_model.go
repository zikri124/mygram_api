package model

import "time"

type StandardClaim struct {
	Jti string `json:"jti"`
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp uint64 `json:"exp"`
	Nbf uint64 `json:"nbf"`
	Iat uint64 `json:"iat"`
}

type AccessClaim struct {
	StandardClaim
	UserID   uint32    `json:"user_id"`
	Username string    `json:"username"`
	DOB      time.Time `json:"dob"`
}
