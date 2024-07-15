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
	getUserFromEmail,err := jsonDatabase.GetUser(userPayload.Email)
	if err != nil {
		sendErrorResponse(w,401,"user not authenticated")
		return
		}
	authenticated := authentication.IsPasswordMatches([]byte(userPayload.Password),[]byte(getUserFromEmail.Password))
	if !authenticated  {
		sendErrorResponse(w,401,"password not matches")
		return
	}
	
	sendJSONResponse(w,getUserFromEmail,200)
}