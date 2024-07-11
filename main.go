package main

import (
	"fmt"
	"net/http"
)


const addr = ":8080"

func main(){
	serveMux := http.ServeMux{}
	apiConfig := apiConfig{}
	apiMux := http.ServeMux{}
	adminMux := http.ServeMux{}
	registerAPIRoutes(&apiMux,&apiConfig)
	registerStaticServerRoutes(&serveMux,&apiConfig)
	registerAdminRoutes(&adminMux,&apiConfig)
	serveMux.Handle("/api/*",&apiMux)
	serveMux.Handle("/admin/*",&adminMux)
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


func registerAPIRoutes(apiMux *http.ServeMux,apiConfig *apiConfig){
	apiMux.HandleFunc("GET /api/healthz",healthCheck)
	apiMux.HandleFunc("GET /api/metrics",apiConfig.getMetrics)
	apiMux.HandleFunc("POST /api/validate_chirp",validateChirpyMessage)
	apiMux.HandleFunc("/api/reset",apiConfig.resetMetrics)
}

func registerAdminRoutes(apiMux *http.ServeMux,apiConfig *apiConfig){
	apiMux.HandleFunc("GET /admin/metrics",apiConfig.getAdminMetrics)
}

func registerStaticServerRoutes(serverMux *http.ServeMux,apiConfig *apiConfig){
	fileServer := http.FileServer(http.Dir("."))
	serverMux.Handle("/app/*",http.StripPrefix("/app",apiConfig.middlewareIncrementHit(fileServer)))
}