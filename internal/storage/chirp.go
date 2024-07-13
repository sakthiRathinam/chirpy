package storage

import (
	"encoding/json"
	"errors"
	"fmt"
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


func (cd *chirpData) AddChirps(chirpMessage string) (bool,error) {
	filePTR,err := os.Create(db_path)
	if err != nil {
		fmt.Println("Error while opening the file")
		return false,errors.New("error while opening the file")
		}
	defer filePTR.Close()
	chirpsData,err := getJsonFileFromStorage(filePTR)
	if err != nil {
		return false,err
	}
	updateChirpsData(&chirpsData,chirpMessage)
	updatedByteData, err := json.Marshal(chirpsData)
	if err != nil {
		return false,err
	}
	err = overwriteJsonToFile(filePTR,updatedByteData)
	if err != nil {
		return false,err
		}
	return true,nil
}


func getJsonFileFromStorage(file *os.File) (databaseStructure,error) {
	fileData := []byte{}
	chirpsData := databaseStructure{}
	_,err :=  file.Read(fileData)
	if err != nil {
		fmt.Println("Error while reading the file")
		return chirpsData,errors.New("error while reading the file")
	}
	json.Unmarshal(fileData,&chirpsData)
	return chirpsData,nil
}

func updateChirpsData(dbStruct *databaseStructure,chirpMessage string) error {
	if dbStruct.Chirp.Chirps == nil {
		dbStruct.Chirp.Chirps = make(map[string]chirp)
	}
	dbStruct.Chirp.IndexCounter = dbStruct.Chirp.IndexCounter + 1
	dbStruct.Chirp.Chirps[fmt.Sprintf("%d",dbStruct.Chirp.IndexCounter)] = chirp{Body:chirpMessage,Id: dbStruct.Chirp.IndexCounter}
	return nil
}

func overwriteJsonToFile(file *os.File,jsonByteData []byte) error {
	file.Write(jsonByteData)
	return nil
} 