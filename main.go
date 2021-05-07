package main

import (
	"github.com/lenalbert/scriptGenYoutube/api"
	"log"
	"net/http"
)

func main() {

	// db := connect

	http.HandleFunc("/generate", api.Generate)

	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Println("Script Generator service failed to start")
		log.Fatal(err)
	}
}
