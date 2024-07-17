package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sakthiRathinam/chirpy/internal/authentication"
)
type user struct {
	Email string `json:"email"`
	Id int `json:"id"`
	Password string `json:"password"`
	RefreshToken string `json:refresh_token`
	token_expiry time.Time
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
	userObj, ok := getUserDataByEmail(userEmail,&dbData)
	if !ok {
		return user{},errors.New("user does not exists")
	}
	return userObj,nil
}


func getUserDataByEmail(userEmail string, dbStruct *databaseStructure) (user,bool){
	for _, userObj := range dbStruct.User.UserData {
		if userObj.Email == userEmail {
			return userObj,true
		}
	}
	return user{},false
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


func (cd *userData) getandUpdateRefreshToken(userEmail string) (user,error) {
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
	userObj, ok := getUserDataByEmail(userEmail,&dbData)
	if !ok {
		return user{},errors.New("user does not exists")
	}
	generateRefreshToken := authentication.CreateRefreshToken()
	userObj.RefreshToken = generateRefreshToken
	userObj.token_expiry = time.Now().Add(time.Duration(1) * time.Hour)
	dbData.User.UserData[fmt.Sprintf("%d",userObj.Id)] = userObj
	return userObj,nil
}


func (cd *userData) updateUser(id,userEmail,password string) (user,error){
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
	userObj,ok := dbData.User.UserData[id]
	if !ok {
		return user{},errors.New("user doesn't exists")
	}
	hashedPassword, _ := authentication.HashPassword(password)
	updatedUserObj := user{Email:userEmail,Id: userObj.Id,Password:hashedPassword}
	dbData.User.UserData[id] = user{Email:userEmail,Id: userObj.Id,Password:hashedPassword}
	updatedByteData, err := json.Marshal(dbData)
	if err != nil {
		return user{},err
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return user{},err
		}
	return updatedUserObj,nil
}


func deleteEmailIfExists(userEmail string,dbStruct *databaseStructure) bool {
	for userKey, userObj := range dbStruct.User.UserData {
		if userObj.Email == userEmail {
			delete(dbStruct.User.UserData,userKey)
			return true
		}
	}
	return false
}

func addUserData(dbStruct *databaseStructure,userEmail string,password string) (user,error) {
	deleteEmailIfExists(userEmail,dbStruct)
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
	dbStruct.User.UserData[fmt.Sprintf("%d",dbStruct.User.IndexCounter)] = userData
	return userData,nil
}