package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log warnings or above
	log.SetLevel(log.WarnLevel)
}

// helloHandler handles HTTP requests to the root URL.
func helloHandler(w http.ResponseWriter, r *http.Request) {
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

	// Use http.Server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,  // Prevent slow-read DoS
		WriteTimeout: 10 * time.Second, // Prevent slow-write DoS
		IdleTimeout:  15 * time.Second, // Limits idle connection reuse
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
