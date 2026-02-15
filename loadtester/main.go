package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		log.Fatal("FATAL: SERVER_HOST environment variable is required.")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		log.Fatal("FATAL: SERVER_PORT environment variable is required.")
	}

	endpoint := "posts"
	targetURL := fmt.Sprintf("http://%s:%s/%s", serverHost, serverPort, endpoint)
	log.Printf("Starting load tester, targeting: %s\n", targetURL)

	// Loop forever to simulate ongoing load
	for {
		resp, err := http.Get(targetURL)
		if err != nil {
			log.Printf("Error connecting to %s: %v\n", targetURL, err)
		} else {
			log.Printf("Request sent! Status: %s\n", resp.Status)
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
	}
}
