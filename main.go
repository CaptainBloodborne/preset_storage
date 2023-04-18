package main

import (
	"fmt"
	"net/http"
	"preset_storage/handlers"
)

func main() {
	fmt.Println("Start server on 8080")
	http.HandleFunc("/get", handlers.HandleGet)
	http.HandleFunc("/set", handlers.HandleSet)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server started!")
}
