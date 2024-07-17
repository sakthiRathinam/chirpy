package main

import (
	"fmt"
	"net/http"

	"github.com/sakthiRathinam/chirpy/internal/authentication"
)



func login(w http.ResponseWriter,r *http.Request){
	userPayload, err := extractPayloadFromUserRequest(r,userRequestPayload{})
	if err != nil {
		sendErrorResponse(w,401,"user not authenticated")
		return
	}
	userObj,err := jsonDatabase.GetUserAndUpdateRefreshToken(userPayload.Email)
	if err != nil {
		sendErrorResponse(w,401,"user not authenticated")
		return
		}
	authenticated := authentication.IsPasswordMatches([]byte(userPayload.Password),[]byte(userObj.Password))
	if !authenticated  {
		sendErrorResponse(w,401,"password not matches")
		return
	}
	generateJWTToken,err := authentication.CreateToken(userObj.Email,userPayload.ExpiresInSeconds,userObj.Id)
	if err != nil {
		sendErrorResponse(w,500,"error while creating token")
		return
	}

	fmt.Println(userObj)
	toSend := map[string]any{"id":userObj.Id,"email":userObj.Email,"token":generateJWTToken,"refresh_token":userObj.RefreshToken}
	sendJSONResponse(w,toSend,200)
}


func refreshAccessToken(w http.ResponseWriter,r *http.Request){
	authHeader := r.Header.Get("Authorization")
	fmt.Println(authHeader,"Header string")
	jwtToken,err := getJWTToken(authHeader)
	if err != nil {
		sendErrorResponse(w,400,"Invalid authorization header")
	}
	userObj,validToken := jsonDatabase.ValidateRefreshToken(jwtToken)
	if !validToken {
		sendErrorResponse(w,401,"invalid token")
		return
	}
	generateJWTToken,err := authentication.CreateToken(userObj.Email,20,userObj.Id)
	if err != nil {
		sendErrorResponse(w,500,"error while creating token")
		return
	}
	toSend := map[string]any{"id":userObj.Id,"email":userObj.Email,"token":generateJWTToken,"refresh_token":userObj.RefreshToken}
	sendJSONResponse(w,toSend,200)
}

func revokeAccessToken(w http.ResponseWriter,r *http.Request){
	authHeader := r.Header.Get("Authorization")
	fmt.Println(authHeader,"Header string")
	jwtToken,err := getJWTToken(authHeader)
	if err != nil {
		sendErrorResponse(w,400,"Invalid authorization header")
	}
	_,revoked := jsonDatabase.RevokeRefreshToken(jwtToken)
	if !revoked {
		sendErrorResponse(w,400,"invalid token")
		return
	}
	sendJSONResponse(w,"revoked",204)
}
