package helper

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

func ValidateToken(token string) (claim jwt.MapClaims, err error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Println("validating jwt error : ", err.Error())
	}

	claim, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Error when translating claim")
		return
	}

	return
}

func GetUserIdFromGinCtx(ctx *gin.Context) (uint32, error) {
	userIdRaw, isExist := ctx.Get("UserId")
	if !isExist {
		return 0, errors.New("cannot get payload in access token")
	}

	userIdFloat := userIdRaw.(float64)
	userId := int(userIdFloat)
	if userId == 0 {
		return 0, errors.New("cannot get payload in access token")
	}

	return uint32(userId), nil
}
