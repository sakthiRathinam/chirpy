package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)


const addr = ":8080"
func main(){
	serveMux := http.ServeMux{}
	serveMux.HandleFunc("/home",homePage)
	fileServer := http.FileServer(http.Dir("."))
	serveMux.Handle("/",fileServer)
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


func homePage(w http.ResponseWriter,r *http.Request){
	htmFile, err := os.Open("index.html")
	if err != nil {
		fmt.Println("error while parsing the request")
		http.Error(w,"Something went wrong",500)
		return 
	}
	defer htmFile.Close()
	w.Header().Set("Content-Type","text/html")

	_,err = io.Copy(w,htmFile)

	if err != nil {
		fmt.Println("error while copying the request")
		http.Error(w,"Something went wrong",500)
		return
	}
}