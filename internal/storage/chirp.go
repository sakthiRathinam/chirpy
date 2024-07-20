package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)
type chirp struct {
	Body string `json:"body"`
	Id int `json:"id"`
	AuthorID int `json:"author_id"`
}

type chirpData struct {
	Chirps map[string]chirp `json:"chirps"`
	IndexCounter int `json:"indexCounter"`
}


func (cd *chirpData) addChirp(chirpMessage string, authorID int) (chirp,error) {
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return chirp{},errors.New("error while opening the file")
		}
	defer filePTR.Close()
	chirpsData,err := getJsonFileFromStorage(filePTR)
	if err != nil {
		return chirp{},err
	}
	chirpObj,err := addChirpsData(&chirpsData,chirpMessage,authorID)
	if err != nil {
		return chirp{},err
	}
	updatedByteData, err := json.Marshal(chirpsData)
	if err != nil {
		return chirp{},err
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return chirp{},err
		}
	return chirpObj,nil
}


func getJsonFileFromStorage(file *os.File) (databaseStructure,error) {
	chirpsData := databaseStructure{}
	fileData,err :=  io.ReadAll(file)
	if err != nil {
		fmt.Println("Error while reading the file")
		return chirpsData,errors.New("error while reading the file")
	}
	json.Unmarshal(fileData,&chirpsData)
	return chirpsData,nil
}

func (cd *chirpData) getAllChirps()([]chirp,error){
	chirpsData := databaseStructure{}
	chirps := []chirp{}
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return chirps,errors.New("error while opening the file")
		}
	defer filePTR.Close()
	fileData,err :=  io.ReadAll(filePTR)
	if err != nil {
		return chirps,err	
	}
	json.Unmarshal(fileData,&chirpsData)
	for _,value := range chirpsData.Chirp.Chirps{
		chirps = append(chirps, value)
	}
	return chirps,nil

}


func (cd *chirpData) getChirp(chirpID string)(chirp,error){
	chirpsData := databaseStructure{}
	chirp := chirp{}
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return chirp,errors.New("error while opening the file")
		}
	defer filePTR.Close()
	fileData,err :=  io.ReadAll(filePTR)
	if err != nil {
		return chirp,err
	}
	json.Unmarshal(fileData,&chirpsData)
	chirp,exists := chirpsData.Chirp.Chirps[chirpID]
	if !exists{
		return chirp, errors.New("chirp not exists")
	}
	return chirp,nil
}

func (cd *chirpData) deleteChirp(chirpID string,authorID int)(bool,error){
	chirpsData := databaseStructure{}
	filePTR,err := os.OpenFile(db_path,os.O_RDWR|os.O_APPEND,7777)
	if err != nil {
		fmt.Println("Error while opening the file")
		return false,errors.New("error while opening the file")
		}
	defer filePTR.Close()
	fileData,err :=  io.ReadAll(filePTR)
	if err != nil {
		return false,nil
	}
	json.Unmarshal(fileData,&chirpsData)
	chirp,exists := chirpsData.Chirp.Chirps[chirpID]
	if !exists{
		return false, errors.New("chirp not exists")
	}
	updatedByteData, err := json.Marshal(chirpsData)
	if err != nil {
		return false,nil
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return false,nil
		}
	if chirp.AuthorID != authorID {
		return false,nil
	}
	fmt.Println(authorID,"author id",chirp.AuthorID,"chirp author id")
	delete(chirpsData.Chirp.Chirps,chirpID)
	return true,nil
}

func addChirpsData(dbStruct *databaseStructure,chirpMessage string,authorID int) (chirp,error) {
	if dbStruct.Chirp.Chirps == nil {
		dbStruct.Chirp.Chirps = make(map[string]chirp)
	}
	dbStruct.Chirp.IndexCounter = dbStruct.Chirp.IndexCounter + 1
	chirpData := chirp{Body:chirpMessage,Id: dbStruct.Chirp.IndexCounter,AuthorID: authorID}
	dbStruct.Chirp.Chirps[fmt.Sprintf("%d",dbStruct.Chirp.IndexCounter)] = chirpData
	return chirpData,nil
}

func overwriteJsonToFile(file *os.File,jsonByteData []byte) error {
	err := file.Truncate(0)
	if err != nil {
		return err
	}
	_,err = file.Seek(0,0)
	if err != nil {
		return err
	}
	file.Write(jsonByteData)
	return nil
} 