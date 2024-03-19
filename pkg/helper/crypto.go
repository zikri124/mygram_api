package helper

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(in string) (out string, err error) {
	outByte, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error when generate hash password : ", err.Error())
		return
	}
	return string(outByte), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(claim any) (token string, err error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	jwtClaim := jwt.MapClaims{}
	encodedClaim, err := json.Marshal(claim)
	if err != nil {
		log.Println("cannot mashal claim payload")
		return
	}

	err = json.Unmarshal(encodedClaim, &jwtClaim)
	if err != nil {
		log.Println("cannot mapping claim to jwt claim")
		return
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaim)

	token, err = parseToken.SignedString([]byte(jwtSecret))

	if err != nil {
		log.Println("cannot generate token : ", err)
		return
	}
	return
}
