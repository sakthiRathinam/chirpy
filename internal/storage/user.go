package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/sakthiRathinam/chirpy/internal/authentication"
)
type user struct {
	Email string `json:"email"`
	Id int `json:"id"`
	Password string `json:"password"` 
}

type userData struct {
	UserData map[string]user `json:"users"`
	IndexCounter int `json:"indexCounter"`
}


func (cd *userData) getUser(userEmail string) (user,error){
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return user{},errors.New("error while opening the file")
		}
	defer filePTR.Close()
	dbData,err := getJsonFileFromStorage(filePTR)
	if err != nil {
		return user{},err
	}
	userObj, ok := dbData.User.UserData[userEmail]
	if !ok {
		return user{},errors.New("user does not exists")
	}
	return userObj,nil
}

func (cd *userData) addUser(userEmail string,password string) (user,error) {
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return user{},errors.New("error while opening the file")
		}
	defer filePTR.Close()
	dbData,err := getJsonFileFromStorage(filePTR)
	if err != nil {
		return user{},err
	}
	userObj,err := addUserData(&dbData,userEmail,password)
	if err != nil {
		return user{},err
	}
	updatedByteData, err := json.Marshal(dbData)
	if err != nil {
		return user{},err
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return user{},err
		}
	return userObj,nil
}

func addUserData(dbStruct *databaseStructure,userEmail string,password string) (user,error) {
	if dbStruct.User.UserData == nil {
		dbStruct.User.UserData = make(map[string]user)
	}
	dbStruct.User.IndexCounter = dbStruct.User.IndexCounter + 1
	hashedPassword, err := authentication.HashPassword(password)
	if err != nil {
		fmt.Println("failed while hashing the password")
		return user{},nil
	}
	userData := user{Email:userEmail,Id: dbStruct.User.IndexCounter,Password: hashedPassword}
	dbStruct.User.UserData[userEmail] = userData
	return userData,nil
}