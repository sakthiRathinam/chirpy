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
	Id int `json:"int"`
}

type chirpData struct {
	Chirps map[string]chirp `json:"chirps"`
	IndexCounter int `json:"indexCounter"`
}


func (cd *chirpData) addChirp(chirpMessage string) (chirp,error) {
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
	chirpObj,err := addChirpsData(&chirpsData,chirpMessage)
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

func addChirpsData(dbStruct *databaseStructure,chirpMessage string) (chirp,error) {
	if dbStruct.Chirp.Chirps == nil {
		dbStruct.Chirp.Chirps = make(map[string]chirp)
	}
	fmt.Println(dbStruct.Chirp)
	dbStruct.Chirp.IndexCounter = dbStruct.Chirp.IndexCounter + 1
	chirpData := chirp{Body:chirpMessage,Id: dbStruct.Chirp.IndexCounter}
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