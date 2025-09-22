package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type RebootRequest struct {
	Token string `json:"token,omitempty"`
	Delay int    `json:"delay,omitempty"` // Delay in seconds
	Force bool   `json:"force,omitempty"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type RebootServer struct {
	authToken string
	mu        sync.Mutex
	rebooting bool
}

func NewRebootServer(authToken string) *RebootServer {
	return &RebootServer{
		authToken: authToken,
	}
}

func (s *RebootServer) authenticate(token string) bool {
	return s.authToken == "" || token == s.authToken
}

func (s *RebootServer) rebootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if already rebooting
	s.mu.Lock()
	if s.rebooting {
		s.mu.Unlock()
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "System is already rebooting",
		})
		return
	}
	s.rebooting = true
	s.mu.Unlock()

	var req RebootRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.mu.Lock()
		s.rebooting = false
		s.mu.Unlock()

		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Authentication
	if !s.authenticate(req.Token) {
		s.mu.Lock()
		s.rebooting = false
		s.mu.Unlock()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Authentication failed",
		})
		return
	}

	// Send immediate response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: fmt.Sprintf("Reboot initiated. Delay: %d seconds", req.Delay),
	})

	// Execute reboot in goroutine (after sending response)
	go s.executeReboot(req.Delay, req.Force)
}

func (s *RebootServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.mu.Lock()
	rebooting := s.rebooting
	s.mu.Unlock()

	status := "ready"
	if rebooting {
		status = "rebooting"
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().Unix(),
	})
}

func (s *RebootServer) executeReboot(delay int, force bool) {
	// Add delay if specified
	if delay > 0 {
		log.Printf("Waiting %d seconds before reboot...", delay)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	log.Println("Initiating system reboot...")

	// Determine reboot command based on force flag
	var cmd *exec.Cmd
	if force {
		cmd = exec.Command("sudo", "reboot", "-f")
	} else {
		cmd = exec.Command("sudo", "reboot")
	}

	// Execute reboot command
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to execute reboot command: %v", err)

		s.mu.Lock()
		s.rebooting = false
		s.mu.Unlock()
		return
	}

	log.Println("Reboot command executed successfully")
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get auth token from environment variable (optional)
	authToken := os.Getenv("REBOOT_AUTH_TOKEN")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewRebootServer(authToken)

	// Set up routes
	http.HandleFunc("/reboot", server.rebootHandler)
	http.HandleFunc("/status", server.statusHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Reboot server starting on port %s", port)
	log.Printf("Authentication: %s", map[bool]string{true: "enabled", false: "disabled"}[authToken != ""])

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
