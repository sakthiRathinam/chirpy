package main

import (
	"fmt"
	"net/http"
)


const addr = ":8080"
func main(){
	serveMux := http.ServeMux{}
	httpServer := http.Server{
		Handler:&serveMux,
		Addr: addr,
	}
	fmt.Printf("Server started on port %s\n",addr)
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("Failed while starting the server",err)
	}
}