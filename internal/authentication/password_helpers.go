package authentication

import "golang.org/x/crypto/bcrypt"




func HashPassword(password string) (string,error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),5)
	if err != nil {
		return "",err
	}
	return string(hashedPassword),nil
}


func IsPasswordMatches(password []byte, hashed_password []byte) bool{
	err := bcrypt.CompareHashAndPassword(hashed_password,password)
	return err == nil
}