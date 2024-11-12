package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
  // Log as JSON instead of the default ASCII formatter
  log.SetFormatter(&log.JSONFormatter{})

  // Output to stdout instead of the default stderr
  // Can be any io.Writer, see below for File example
  log.SetOutput(os.Stdout)

  // Only log the warning severity or above.
  log.SetLevel(log.WarnLevel)
}

// hello world from 1
// helloHandler handles HTTP requests to the root URL.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// This will return the string without.
	fmt.Fprint(w, "Hello, World!") 
}

// main function sets up the HTTP server.
func main() {
	log.WithFields(log.Fields{
    "animal": "walrus",
    "size":   10,
  }).Info("A group of walrus emerges from the ocean")

	http.HandleFunc("/", helloHandler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
