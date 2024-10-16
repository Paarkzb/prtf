package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	msg := os.Getenv("MESSAGE")

	http.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %s", msg)
	})

	log.Printf("Starting service %s", msg)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
