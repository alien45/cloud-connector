package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var err error
var err400 = http.StatusBadRequest
var err500 = http.StatusInternalServerError
var err404 = http.NotFound
var ok200 = http.StatusOK
var created201 = http.StatusCreated

func main() {
	args := os.Args[1:]
	port := "3000"
	if len(args) >= 1 {
		port = args[0]
	}

	if len(args) >= 2 && strings.ToLower(args[1]) == "demo" {
		// Serve ./public directory as web page for testing file upload
		http.Handle("/", http.FileServer(http.Dir("./public")))
		log.Printf("Visit http://localhost:%s to test file upload\n", port)
	}
	// Register handers for API endpoints
	http.HandleFunc("/dial", dialHandler)
	http.HandleFunc("/containers", containersHandler)
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/copy", copyHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/getjsonstring", getJSONStringHandler)

	log.Println("Cloud connector started at port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func respondJSON(w http.ResponseWriter, content interface{}, statusCode int) {
	b, err := json.Marshal(content)
	if respondIfError(err, w, "Something went wrong!", err500) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(b)
	if err != nil {
		log.Println("[response] [error]", err)
		return
	}

	log.Printf("[response] [status%d] %s\n", statusCode, http.StatusText(statusCode))
}

func respondIfError(err error, w http.ResponseWriter, msg string, statusCode int) bool {
	if err == nil {
		return false
	}
	respondError(w, msg, statusCode)
	return true
}

func respondError(w http.ResponseWriter, msg string, statusCode int) {
	if statusCode == 0 {
		statusCode = err400
	}
	if msg == "" {
		msg = http.StatusText(statusCode)
	}
	http.Error(w, msg, statusCode)
	log.Printf("[response] [status%d] %s\n", statusCode, msg)
}
