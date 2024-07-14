package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)
type chirpRequestPayload struct {
	Body string `json:"body"`
}


func validateChirpyMessage(w http.ResponseWriter,r *http.Request){
	chirpMessage,err := extractPayloadFromChirpRequest(r,chirpRequestPayload{})
	if err != nil {
		sendErrorResponse(w,500,"")
	}

	length_of_message := len(chirpMessage.Body)
	if length_of_message > 140 {
		sendErrorResponse(w,400,"Chirp is too long")
		return
	}

	data := map[string]string{"valid":"true","cleaned_body":replaceProfaneWords(chirpMessage.Body)}
	sendJSONResponse(w,data,201)
	
}

func addChirp(w http.ResponseWriter,r *http.Request){
	chirpMessage, err := extractPayloadFromChirpRequest(r,chirpRequestPayload{})
	if err != nil {
		sendErrorResponse(w,500,"")
		return
		}
	length_of_message := len(chirpMessage.Body)
	if length_of_message > 140 {
		sendErrorResponse(w,400,"Chirp is too long")
		return
	}
	chirp,err :=jsonDatabase.AddChirp(chirpMessage.Body)
	if err != nil {
		fmt.Println("Failed while adding chirp")
		sendErrorResponse(w,500,"failed while adding the chirp, please try again!!!!")
		}
	sendJSONResponse(w,chirp,201)
}

func getAllChirps(w http.ResponseWriter,r *http.Request){
	chirpObjs,err := jsonDatabase.GetChirps()
	if err != nil{
		sendErrorResponse(w,500,"error while fetching the chirps")
		return
	}
	sendJSONResponse(w,chirpObjs,200)
}


func getChirp(w http.ResponseWriter,r *http.Request){
	chirpID := r.PathValue("chirpID")
	fmt.Println(chirpID)
	if chirpID == ""{
		sendErrorResponse(w,400,"path value is not valid")
		return
	}

	chirp,err := jsonDatabase.GetChirp(chirpID)
	if err != nil {
		sendErrorResponse(w,404,"chirp id not found")
		return
	}
	sendJSONResponse(w,chirp,200)
}


func extractPayloadFromChirpRequest(r *http.Request,payload chirpRequestPayload) (chirpRequestPayload,error) {
	requestBody := payload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		fmt.Println("Error decoding paramerts",err)
		return requestBody,err
	}
	return requestBody,nil
}


func replaceProfaneWords(body string) string {
	profaneWords := []string{"kerfuffle","sharbert","fornax"}
	for _, profanceWord := range profaneWords{
		body = replaceCaseInsensitive(body,profanceWord,"****")
	}
	return body
}


func replaceCaseInsensitive(input, old,new string)string {
	re, err := regexp.Compile("(?i)" + old)
	if err != nil {
		fmt.Println("Error compiling regex:",err)
		return input
	}

	result := re.ReplaceAllString(input, new)
	return result
}