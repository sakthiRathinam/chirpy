package main

import (
	"net/http"

	"github.com/sakthiRathinam/chirpy/internal/authentication"
)



func login(w http.ResponseWriter,r *http.Request){
	userPayload, err := extractPayloadFromUserRequest(r,userRequestPayload{})
	if err != nil {
		sendErrorResponse(w,401,"user not authenticated")
		return
	}
	userObj,err := jsonDatabase.GetUser(userPayload.Email)
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


	toSend := map[string]any{"id":userObj.Id,"email":userObj.Email,"token":generateJWTToken}
	sendJSONResponse(w,toSend,200)
}