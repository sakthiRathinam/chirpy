package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type userRequestPayload struct {
	Email string `json:"email"`
}
func addUser(w http.ResponseWriter,r *http.Request){
	userMail, err := extractPayloadFromUserRequest(r,userRequestPayload{})
	if err != nil {
		sendErrorResponse(w,500,"")
	}
	fmt.Println(userMail)
	fmt.Printf("%T",userMail)
	user,err :=jsonDatabase.AddUser(userMail.Email)
	if err != nil {
		fmt.Println("Failed during adding new user")
		sendErrorResponse(w,500,"failed while adding the user, please try again!!!!")
	}
	sendJSONResponse(w,user,201)
}


func extractPayloadFromUserRequest(r *http.Request,payload userRequestPayload) (userRequestPayload,error) {
	requestBody := payload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		fmt.Println("Error decoding paramerts",err)
		return requestBody,err
	}
	return requestBody,nil
}