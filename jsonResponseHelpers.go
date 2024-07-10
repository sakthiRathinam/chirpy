package main

import (
	"encoding/json"
	"net/http"
)


func sendJSONResponse(w http.ResponseWriter, toSend interface{},statusCode int) error {
	data, err := json.Marshal(toSend)
	if err != nil {
		return sendErrorResponse(w,500)
	}
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
	return nil

}


func sendErrorResponse(w http.ResponseWriter, statusCode int) error {
	return sendJSONResponse(w,map[string]string{"Error":"Something went wrong"},statusCode)
}