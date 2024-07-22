package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/sakthiRathinam/chirpy/internal/authentication"
)

type userRequestPayload struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func addUser(w http.ResponseWriter, r *http.Request) {
	userPayload, err := extractPayloadFromUserRequest(r, userRequestPayload{})
	if err != nil {
		sendErrorResponse(w, 500, "")
	}
	user, err := jsonDatabase.AddUser(userPayload.Email, userPayload.Password)
	if err != nil {
		fmt.Println("Failed during adding new user")
		sendErrorResponse(w, 500, "failed while adding the user, please try again!!!!")
	}
	sendJSONResponse(w, user, 201)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userPayload, err := extractPayloadFromUserRequest(r, userRequestPayload{})
	if err != nil {
		sendErrorResponse(w, 500, "")
	}

	authHeader := r.Header.Get("Authorization")
	jwtToken, err := getJWTToken(authHeader)

	if err != nil {
		sendErrorResponse(w, 401, "invalid token")
		return
	}

	id, err := authentication.ValidateAndExtractIDFromToken(jwtToken)
	if err != nil {
		sendErrorResponse(w, 401, "invalid token")
		return
	}
	updatedUser, err := jsonDatabase.UpdateUser(id, userPayload.Email, userPayload.Password)
	if err != nil {
		sendErrorResponse(w, 401, "invalid token")
		return
	}
	sendJSONResponse(w, updatedUser, 200)
}

func getJWTToken(token string) (string, error) {
	splittedArr := strings.Split(token, " ")
	if len(splittedArr) != 2 {
		return "", errors.New("invalid token")
	}
	return splittedArr[1], nil
}

func extractPayloadFromUserRequest(r *http.Request, payload userRequestPayload) (userRequestPayload, error) {
	requestBody := payload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		return requestBody, err
	}
	return requestBody, nil
}
