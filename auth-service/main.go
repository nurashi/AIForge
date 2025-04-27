package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "handle func in Main of Auth service")
	})

	log.Println("server starts on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("ERROR ERROR ERROR: ", err)

	}
}
