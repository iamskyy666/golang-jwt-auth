package main

import (
	"log"
	"net/http"
	"time"

	"github.com/callmeskyy111/golang-jwt-auth/internal/httpserver"
)

func main() {

	router := httpserver.NewRouter()

	// std. golang way to run a http-server.
	srv:=&http.Server{
		Addr:":5000",
		Handler: router,
		ReadHeaderTimeout: time.Second * 5,
	}

	log.Printf("API running on PORT %s ✅",srv.Addr)

	if err:=srv.ListenAndServe(); err!=nil{
		if err==http.ErrServerClosed{
			log.Printf("Server Closed!")
			return
		}
		log.Fatalf("⚠️ Server Error: %v",err)
	}

}