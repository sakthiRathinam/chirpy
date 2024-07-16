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
	if expires_at == 0 {
		expires_at = 2
	}
	claims := ChirpyCliam{
		email,
		jwt.RegisteredClaims{Issuer:"chirpy",
		Subject:"authenticate chirpy user",           
		Audience:[]string{"chipryuser","chirp"},  
		ExpiresAt:jwt.NewNumericDate(time.Now().Add(time.Duration(expires_at) * time.Minute)),
		NotBefore:jwt.NewNumericDate(time.Now()),
		IssuedAt:jwt.NewNumericDate(time.Now()),   
		ID:fmt.Sprintf("%d",id)}, 
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signedToken,err := token.SignedString([]byte(signedKey))
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed while creating the token")
		return "",err
	}
	return signedToken,nil
}

func ValidateAndExtractIDFromToken(token string) (string,error) {
	claim := ChirpyCliam{}
	jwtToken,err := jwt.ParseWithClaims(token,&claim,func(t *jwt.Token) (interface{},error) {
		signedKeyInBytes := []byte(signedKey)
		return signedKeyInBytes, nil
	})
	if  err!= nil{
		return "",errors.New("invalid token")
	}
	if !jwtToken.Valid{
		return "",errors.New("invalid token")
	}
	fmt.Println(claim.ID,"id hereeeee",token)
	return claim.ID,nil
}