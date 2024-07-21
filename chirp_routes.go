package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/sakthiRathinam/chirpy/internal/authentication"
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
	authHeader := r.Header.Get("Authorization")
	jwtToken,err := getJWTToken(authHeader)
	
	if err != nil {
		sendErrorResponse(w,401,"invalid token")
		return
	}

	stringID, err := authentication.ValidateAndExtractIDFromToken(jwtToken)
	if err != nil {
		sendErrorResponse(w,401,"invalid token")
		return
	}
	integerID,err := strconv.Atoi(stringID)
	if err != nil {
        fmt.Println("Error:", err)
        return
    }
	chirp,err :=jsonDatabase.AddChirp(chirpMessage.Body,integerID)
	if err != nil {
		fmt.Println("Failed while adding chirp")
		sendErrorResponse(w,500,"failed while adding the chirp, please try again!!!!")
		}
	sendJSONResponse(w,chirp,201)
}

func getAllChirps(w http.ResponseWriter,r *http.Request){
	authorID := r.URL.Query().Get("author_id")
	if authorID != ""{
		integerAuthID,err := strconv.Atoi(authorID)
		if err != nil {
			sendErrorResponse(w,400,"Invalid author id format")
			return
		}
		chirpObjs,err := jsonDatabase.GetChirpsForAuthorID(integerAuthID)
		if err != nil{
		sendErrorResponse(w,500,"error while fetching the chirps")
		return
		}
		sendJSONResponse(w,chirpObjs,200)
		return
	}
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


func deleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	authHeader := r.Header.Get("Authorization")
	jwtToken,err := getJWTToken(authHeader)
	
	if err != nil {
		sendErrorResponse(w,401,"invalid token")
		return
	}

	stringID, err := authentication.ValidateAndExtractIDFromToken(jwtToken)
	if err != nil {
		sendErrorResponse(w,401,"invalid token")
		return
	}
	integerID,err := strconv.Atoi(stringID)
	if err != nil {
        fmt.Println("Error:", err)
        return
    }
	fmt.Println(chirpID,integerID,"from the response")
	deleteChirp,err := jsonDatabase.DeleteChirp(chirpID,integerID)
	if err != nil {
		sendErrorResponse(w,400,"Not authorized")
		return
	}
	if !deleteChirp {
		sendErrorResponse(w,403,"Not authorized")
		return
	}
	sendJSONResponse(w,"deleted",204)

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