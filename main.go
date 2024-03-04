package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var logger *log.Logger

func init() {
	// Open or create a log file
	file, err := os.OpenFile("proxy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	// Set log output to the file
	logger = log.New(file, "Proxy: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	// logger = log.New(file, "Proxy: ", log.Ldate|log.Ltime|log.Lshortfile)
	// Also log to the console for visibility during development
	logger.SetOutput(io.MultiWriter(file, os.Stdout))
}

func proxyGetHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Extract query parameters from the request
	appName := r.URL.Query().Get("app_name")
	appPort := r.URL.Query().Get("app_port")
	appEndpoint := r.URL.Query().Get("app_endpoint")

	// Construct the target URL
	targetURL := fmt.Sprintf("http://%s:%s/%s", appName, appPort, appEndpoint)
	logger.Println(targetURL)

	// Perform the proxy request
	response, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Error fetching metrics: "+err.Error(), http.StatusInternalServerError)
		logger.Println("Proxy request failed:", err)
		return
	}
	defer response.Body.Close()

	// Copy the response from the target server to the original requester
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.WriteHeader(response.StatusCode)
	_, err = io.Copy(w, response.Body)
	if err != nil {
		logger.Println("Error copying response:", err)
	}

	duration := time.Since(startTime)
	logger.Printf("Proxy request duration: %v", duration)
	logger.Printf("Proxy request status code: %d", response.StatusCode)
}

func proxyPostHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Extract query parameters from the request
	appName := r.URL.Query().Get("app_name")
	appPort := r.URL.Query().Get("app_port")
	appEndpoint := r.URL.Query().Get("app_endpoint")

	// Parse JSON body if present
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil && err != io.EOF {
		http.Error(w, "Error decoding JSON body: "+err.Error(), http.StatusBadRequest)
		logger.Println("Proxy request failed:", err)
		return
	}

	// Construct the target URL
	targetURL := fmt.Sprintf("http://%s:%s/%s", appName, appPort, appEndpoint)
	logger.Println(targetURL)

	// Perform the proxy request with JSON body
	client := &http.Client{}
	req, err := http.NewRequest("POST", targetURL, nil)
	if err != nil {
		http.Error(w, "Error creating POST request: "+err.Error(), http.StatusInternalServerError)
		logger.Println("Proxy request failed:", err)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// If JSON body is present, set it in the request
	if len(requestBody) > 0 {
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			http.Error(w, "Error encoding JSON body: "+err.Error(), http.StatusInternalServerError)
			logger.Println("Proxy request failed:", err)
			return
		}
		req.Body = io.NopCloser(bytes.NewReader(jsonBody))
	}

	// Perform the proxy request
	response, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error fetching metrics: "+err.Error(), http.StatusInternalServerError)
		logger.Println("Proxy request failed:", err)
		return
	}
	defer response.Body.Close()

	// Copy the response from the target server to the original requester
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.WriteHeader(response.StatusCode)
	_, err = io.Copy(w, response.Body)
	if err != nil {
		logger.Println("Error copying response:", err)
	}

	duration := time.Since(startTime)
	logger.Printf("Proxy request duration: %v", duration)
	logger.Printf("Proxy request status code: %d", response.StatusCode)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/proxyGet/", proxyGetHandler).Queries("app_name", "{app_name}", "app_port", "{app_port}", "app_endpoint", "{app_endpoint}").Methods("GET")
	router.HandleFunc("/proxyPost/", proxyPostHandler).Queries("app_name", "{app_name}", "app_port", "{app_port}", "app_endpoint", "{app_endpoint}").Methods("POST")

	// Start the server
	log.Println("Server listening on :5000")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
