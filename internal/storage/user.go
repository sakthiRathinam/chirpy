package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/sakthiRathinam/chirpy/internal/authentication"
)
type user struct {
	Email string `json:"email"`
	Id int `json:"id"`
	Password string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	TokenExpiry time.Time `json:"token_expiry"`
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
	userObj, ok := getUserDataByStructAttr(userEmail,"Email",&dbData)
	if !ok {
		return user{},errors.New("user does not exists")
	}
	generateRefreshToken := authentication.CreateRefreshToken()
	dbData.User.UserData[fmt.Sprintf("%d",userObj.Id)] = user{Email:userEmail,Id: userObj.Id,Password:userObj.Password,RefreshToken: generateRefreshToken,TokenExpiry: time.Now().Add(time.Duration(1) * time.Hour)}
	fmt.Println("updated token",generateRefreshToken)
	updatedByteData, err := json.Marshal(dbData)
	if err != nil {
		return user{},err
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return user{},err
		}
	return dbData.User.UserData[fmt.Sprintf("%d",userObj.Id)],nil
}


func (cd *userData) validateRefreshToken(refershToken string) (user,bool) {
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return user{},false
		}
	defer filePTR.Close()
	dbData,err := getJsonFileFromStorage(filePTR)
	if err != nil {
		return user{},false
	}
	userObj, ok := getUserDataByStructAttr(refershToken,"RefreshToken",&dbData)
	fmt.Println(userObj,ok,"refressh toke founddddddddddddddddddddddddd")
	if !ok {
		return user{},false
	}

	if userObj.TokenExpiry.Before(time.Now()) {
		return userObj,false
	}
	return userObj,true
}

func (cd *userData) revokeRefreshToken(refershToken string) (user,bool) {
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return user{},false
		}
	defer filePTR.Close()
	dbData,err := getJsonFileFromStorage(filePTR)
	if err != nil {
		return user{},false
	}
	userObj, ok := getUserDataByStructAttr(refershToken,"RefreshToken",&dbData)
	if !ok {
		return user{},false
	}

	fmt.Println("user token got revoked",userObj.RefreshToken)
	userObj.TokenExpiry = time.Now().Add(-time.Duration(5 * time.Second))
	dbData.User.UserData[fmt.Sprintf("%d",userObj.Id)] = userObj
	updatedByteData, err := json.Marshal(dbData)
	if err != nil {
		return user{},false
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return user{},false
		}
	return userObj,true

}


func getUserDataByStructAttr(attr string,attrName string,dbStruct *databaseStructure) (user,bool){
	for _, userObj := range dbStruct.User.UserData {
		value , err := getAttribute(userObj,attrName)
		if err != nil {
			continue
		}
		if value == attr {
			return userObj,true
		}
	}
	return user{},false
}


func getAttribute(obj interface{}, attrName string) (interface{}, error) {
	v := reflect.ValueOf(obj)
	field := v.FieldByName(attrName)
	if !field.IsValid() {
		return nil, fmt.Errorf("no such field: %s", attrName)
	}
	return field.Interface(), nil
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
	fmt.Println(userObj.RefreshToken,"before updating")
	if !ok {
		return user{},errors.New("user doesn't exists")
	}
	hashedPassword, _ := authentication.HashPassword(password)
	dbData.User.UserData[id] = user{Email:userEmail,Id: userObj.Id,Password:hashedPassword,RefreshToken: userObj.RefreshToken,TokenExpiry: userObj.TokenExpiry}
	updatedByteData, err := json.Marshal(dbData)
	if err != nil {
		return user{},err
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return user{},err
		}
	return dbData.User.UserData[id],nil
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