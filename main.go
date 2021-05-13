package main

import (
	"./api"
	"log"
	"net/http"
)

func main() {

	log.Println("Script Generator Service Started !")
	// db := connect

	http.HandleFunc("/generate", api.Generate)

	err := http.ListenAndServe(":10001", nil)
	if err != nil {
		log.Println("Script Generator service failed to start")
		log.Fatal(err)
	}
}
