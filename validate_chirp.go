package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)
type requestPayload struct {
	Body string `json:"body"`
}


func validateChirpyMessage(w http.ResponseWriter,r *http.Request){
	requestBody := requestPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		fmt.Println("Error decoding paramerts",err)
		sendErrorResponse(w,500,"")
		return 
	}

	length_of_message := len(requestBody.Body)

	if length_of_message > 140 {
		sendErrorResponse(w,400,"Chirp is too long")
		return
	}
	data, err := json.Marshal(map[string]string{"valid":"true","cleaned_body":replaceProfaneWords(requestBody.Body)})
	if err != nil {
		fmt.Println("Error decoding paramerts",err)
		sendErrorResponse(w,500,"")
		return 
	}
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
	w.Write(data) 
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