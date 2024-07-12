package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)




const db_path = "database.json"

type chirp struct {
	Body string `json:"body"`
	Id int `json:"int"`
}

type chirpData struct {
	Chirps map[string]chirp `json:"chirps"`
	IndexCounter int `json:"indexCounter"`
}


func (cd *chirpData) AddChirps(chirpMessage string) (bool,error) {
	filePTR,err := os.Create("database.json")
	if err != nil {
		fmt.Println("Error while opening the file")
		return false,errors.New("error while opening the file")
	}
	byteData := []byte{}
	chirpsData := databaseStructure{}
	lenOfData,err :=  filePTR.Read(byteData)
	if err != nil {
		fmt.Println("Error while reading the file")
		return false,nil
	}
	json.Unmarshal(byteData,&chirpsData)
	if chirpsData.Chirp.Chirps == nil {
		chirpsData.Chirp.Chirps = make(map[string]chirp)
	}
	fmt.Println(chirpsData.Chirp.IndexCounter)
	chirpsData.Chirp.IndexCounter = chirpsData.Chirp.IndexCounter + 1
	chirpsData.Chirp.Chirps[fmt.Sprintf("%d",chirpsData.Chirp.IndexCounter)] = chirp{Body:chirpMessage,Id: chirpsData.Chirp.IndexCounter}
	byteData, err = json.Marshal(chirpsData)
	if err != nil {
		fmt.Println("Error while parsing the bytes")
	}
	filePTR.Write(byteData)
	defer filePTR.Close()
	fmt.Println(chirpsData,lenOfData)
	return true,nil
}

func overwriteJsonToFile(jsonByteData []byte) error {
	filePTR,err := os.Create("database.json")
	if err != nil {
		fmt.Println("Error while writing the bytes")
		return errors.New("Error while writing the bytes")
	}
	filePTR.Write(jsonByteData)
	defer filePTR.Close()
	return nil
} 

type databaseStructure struct {
	Chirp chirpData `json:"chirp"`
	
}

type JsonDatabase struct {
	DB databaseStructure
	RMtx *sync.RWMutex
}


func (jd *JsonDatabase) EnsureDB() error {
	fileExists := fileExists(db_path)
	if !fileExists {
		err := createDatabaseFile(db_path)
		if err != nil{
			panic(err)
		}
	}
	return nil
}


func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


func createDatabaseFile(path string) error {
	filePTR,_ := os.Create(path)
	dummyStructure := databaseStructure{
			Chirp:chirpData{Chirps:map[string]chirp{},
			IndexCounter:0},
		}
	jsonBytes,err := json.Marshal(dummyStructure)
	filePTR.Write(jsonBytes)
	defer filePTR.Close()
	if err != nil {
		fmt.Println("failed during file creation")
	}
	return nil
}

func appendDummyStructre() error {
	return nil
}

func CreateJsonDatabase() *JsonDatabase {
	jsonDatabse := JsonDatabase{RMtx:&sync.RWMutex{}}
	return &jsonDatabse
}