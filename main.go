package main

import (
	"log"
	"net/http"

	"github.com/go-services/gomod/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err2 := godotenv.Load()
	if err2 != nil {
		log.Fatal("Not able to load environemnet variable")
	}
	Router := mux.NewRouter()
	// Or extend your config for customization

	Router.HandleFunc("/upload", services.Upload).Methods("PUT")
	Router.HandleFunc("/download", services.Download).Methods("GET")
	log.Println("Server listening at port 3000")
	http.Handle("/", Router)
	http.ListenAndServe(":3000", nil)

}
