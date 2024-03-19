package helper

import (
	"log"

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
