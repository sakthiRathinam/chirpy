package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)


type WebhookEvent struct {
	Event string `json:"event"`
	Data struct{ 
		UserID int `json:"user_id"` 
	} `json:"data"`
}
func polkaWebhook(w http.ResponseWriter, r *http.Request) {
	eventPayload,err := extractPayloadFromWebhookEvent(r)
	if err != nil {
		sendErrorResponse(w,400,"invalid format")
		return
	}
	authHeader := r.Header.Get("Authorization")
	apiKey,err := extractApiKeyFromReq(authHeader)
	
	if err != nil {
		sendErrorResponse(w,401,"invalid token")
		return
	}
	polkaApiKey := os.Getenv("POLKA_API_KEY")
	fmt.Println(polkaApiKey,apiKey)
	if apiKey != polkaApiKey {
		sendErrorResponse(w,402,"api key not valid")
		return
	}
	fmt.Println("received webhook event",eventPayload)
	if eventPayload.Event == "user.upgraded" {
		_,err := jsonDatabase.UpdateChirpyRed(eventPayload.Data.UserID)
		if err != nil {
			sendErrorResponse(w,404,"user not found")
			return
		}
		fmt.Println("user got updated")
		sendJSONResponse(w,"user updated",204)
		return
		}
	sendJSONResponse(w,"someother event",204)
}


func extractApiKeyFromReq(authHeader string) (string,error){
	splittedArr := strings.Split(authHeader, " ")
	if len(splittedArr) != 2 {
		return "",errors.New("invalid token")
	}
	return splittedArr[1],nil
}
func extractPayloadFromWebhookEvent(r *http.Request) (WebhookEvent,error) {
	requestBody := WebhookEvent{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		return requestBody,err
	}
	return requestBody,nil
}