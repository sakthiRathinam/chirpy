package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)
type user struct {
	Email string `json:"email"`
	Id int `json:"id"`
	password string 
}

type userData struct {
	UserData map[string]user `json:"chirps"`
	IndexCounter int `json:"indexCounter"`
}


func (cd *userData) addUser(userEmail string) (user,error) {
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return user{},errors.New("error while opening the file")
		}
	defer filePTR.Close()
	chirpsData,err := getJsonFileFromStorage(filePTR)
	if err != nil {
		return user{},err
	}
	userObj,err := addUserData(&chirpsData,userEmail)
	if err != nil {
		return user{},err
	}
	updatedByteData, err := json.Marshal(chirpsData)
	if err != nil {
		return user{},err
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return user{},err
		}
	return userObj,nil
}


func addUserData(dbStruct *databaseStructure,userEmail string) (user,error) {
	if dbStruct.User.UserData == nil {
		dbStruct.User.UserData = make(map[string]user)
	}
	dbStruct.Chirp.IndexCounter = dbStruct.Chirp.IndexCounter + 1
	userData := user{Email:userEmail,Id: dbStruct.Chirp.IndexCounter}
	dbStruct.User.UserData[userEmail] = userData
	return userData,nil
}