package authentication

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signedKey = os.Getenv("JWT_SECRET")
type ChirpyCliam struct {
	Email string 
	jwt.RegisteredClaims
}


func CreateToken(email string, expires_at int, id int) (string,error){
	claims := ChirpyCliam{
		email,
		jwt.RegisteredClaims{Issuer:"chirpy",
		Subject:"authenticate chirpy user",           
		Audience:[]string{"chipryuser","chirp"},  
		ExpiresAt:jwt.NewNumericDate(time.Now().Add(time.Duration(expires_at) * time.Second)),
		NotBefore:jwt.NewNumericDate(time.Now()),
		IssuedAt:jwt.NewNumericDate(time.Now()),   
		ID:fmt.Sprintf("%d",id)}, 
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signedToken,err := token.SignedString([]byte(signedKey))
	if err != nil {
		fmt.Println("failed while creating the token")
		return "",err
	}
	return signedToken,nil
}

func ValidateToken(token string) (string,error) {
	claim := ChirpyCliam{}
	_,err := jwt.ParseWithClaims(token,&claim,func(t *jwt.Token) (interface{},error) {
		if !t.Valid{
			return "",errors.New("invalid token")
		}
		return "",nil
	})
	if err != nil {
		return "",errors.New("invalid token")
	}
	return claim.ID,nil
}