package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)




const db_path = "database.json"


type databaseStructure struct {
	Chirp chirpData `json:"chirp"`
	User userData `json:"user"`
	
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

func (jd *JsonDatabase) FlushDB() error {
	fileExists := fileExists(db_path)
	if fileExists {
		err := os.Remove(db_path)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (jd *JsonDatabase) AddChirp(chirpMessage string) (chirp,error) {
	jd.RMtx.Lock()
	chirp,err := jd.DB.Chirp.addChirp(chirpMessage)
	defer jd.RMtx.Unlock()
	return chirp,err
}

func (jd *JsonDatabase) AddUser(userEmail string,password string) (user,error) {
	jd.RMtx.Lock()
	userObj,err := jd.DB.User.addUser(userEmail,password)
	defer jd.RMtx.Unlock()
	return userObj,err
}
func (jd *JsonDatabase) UpdateUser(id,userEmail,password string) (user,error) {
	jd.RMtx.Lock()
	userObj,err := jd.DB.User.updateUser(id,userEmail,password)
	defer jd.RMtx.Unlock()
	return userObj,err
}

func (jd *JsonDatabase) GetUser(userEmail string) (user,error) {
	jd.RMtx.Lock()
	userObj,err := jd.DB.User.getUser(userEmail)
	defer jd.RMtx.Unlock()
	return userObj,err
}

func (jd *JsonDatabase) GetChirps() ([]chirp,error) {
	jd.RMtx.RLock()
	chirp,err := jd.DB.Chirp.getAllChirps()
	defer jd.RMtx.RUnlock()
	return chirp,err
}

func (jd *JsonDatabase) GetChirp(chirpID string) (chirp,error) {
	jd.RMtx.RLock()
	chirp,err := jd.DB.Chirp.getChirp(chirpID)
	defer jd.RMtx.RUnlock()
	return chirp,err
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

func CreateJsonDatabase() *JsonDatabase {
	jsonDatabse := JsonDatabase{RMtx:&sync.RWMutex{}}
	return &jsonDatabse
}