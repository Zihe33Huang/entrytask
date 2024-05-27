package tcp

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var secretKey = []byte("zihehuang")
var Salt = "zihehuang"

//func Encode(password string) (string, error) {
//	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	return string(hash), err
//}

//
//func compare(hashedPassword, password string) error {
//	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
//	return err
//}

func Encode(password string) (string, error) {
	// Hash password using MD5
	hashedPassword := md5.Sum([]byte(password + Salt))
	return hex.EncodeToString(hashedPassword[:]), nil
}

func compare(hashedPassword, password string) error {
	// Hash the incoming password and compare it with the stored hash
	hashedInput := md5.Sum([]byte(password + Salt))
	hashedInputStr := hex.EncodeToString(hashedInput[:])

	if hashedInputStr != hashedPassword {
		return fmt.Errorf("Passwords Do Not Match")
	}

	return nil
}

func generateJwt(userId uint64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix() // Token expiration time (1 day)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	// Return the token as a response
	return tokenString, err
}
