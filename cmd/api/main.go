package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/callmeskyy111/golang-jwt-auth/internal/app"
	"github.com/callmeskyy111/golang-jwt-auth/internal/httpserver"
)

// main-entry file
func main() {

	ctx:=context.Background() // root ctxt.

	a,err:=app.NewApp(ctx)
	if err != nil {
		log.Fatalf("⚠️ App-Startup failed: %v",err)
	}
	defer func(){
		if err:=a.CloseMongo(ctx); err!=nil{
			log.Printf("⚠️ ShutDown warning: %v",err)
		}
	}()


	router := httpserver.NewRouter(a)

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
