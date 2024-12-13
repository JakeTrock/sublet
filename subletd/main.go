package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// Message represents the WebSocket message structure
type Message struct {
	Type      string            `json:"type"`
	Message   string            `json:"message"`
	Command   string            `json:"command,omitempty"`
	Files     map[string]string `json:"files,omitempty"`
	Timestamp string            `json:"timestamp,omitempty"`
	ID        string            `json:"id,omitempty"`
	TargetID  string            `json:"targetId,omitempty"`
}

func handleCapabilityCommand(command string, files map[string]string) (string, error) {
	switch command {
	case "init":
		err := initConfiguration()
		if err != nil {
			return "", err
		}
		return "Configuration initialized", nil

	case "list":
		files, err := listNixFiles()
		if err != nil {
			return "", err
		}
		return "Files: " + strings.Join(files, ", "), nil

	case "get":
		contents, err := getNixFilesContents()
		if err != nil {
			return "", err
		}
		return contents, nil

	case "set":
		if files == nil {
			return "", fmt.Errorf("no files provided")
		}
		err := setNixFilesContents(files)
		if err != nil {
			return "", err
		}
		return "Files updated successfully", nil

	case "dry-build":
		output, err := runDryBuild()
		if err != nil {
			return "", err
		}
		return output, nil

	case "test":
		output, err := runTest()
		if err != nil {
			return "", err
		}
		return output, nil

	case "switch":
		output, err := runSwitch()
		if err != nil {
			return "", err
		}
		return output, nil

	case "liveness":
		output, err := livenessCheck()
		if err != nil {
			return "", err
		}
		return output, nil

	case "shutdown":
		output, err := shutdown()
		if err != nil {
			return "", err
		}
		return output, nil

	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
}

// Generate unique ID
// func generateID() string {
// 	bytes := make([]byte, 16)
// 	if _, err := rand.Read(bytes); err != nil {
// 		log.Fatal(err)
// 	}
// 	return hex.EncodeToString(bytes)
// }

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Error: Usage: subletd <host_url> <client_id>")
	}

	hostURL := strings.TrimSpace(os.Args[1])
	if hostURL == "" {
		log.Fatal("Error: Host URL cannot be empty")
	}

	clientID := strings.TrimSpace(os.Args[2])
	if clientID == "" {
		log.Fatal("Error: Client ID cannot be empty")
	}
	// ensure clientID is a valid UUID(length of 32)
	if len(clientID) != 32 {
		log.Fatal("Error: Client ID must be a valid UUID")
	}

	log.Printf("Using host URL: %s", hostURL)
	log.Printf("Using client ID: %s", clientID)

	// Connect to WebSocket server
	u := url.URL{Scheme: "ws", Host: hostURL, Path: "/ws"}
	ws, err := websocket.Dial(u.String(), "", hostURL)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer ws.Close()

	// Set up channels
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	done := make(chan struct{})

	// Wait for initial system message
	var initialMsg string
	if err := websocket.Message.Receive(ws, &initialMsg); err != nil {
		log.Fatal("Initial message error:", err)
	}
	log.Printf("Received initial message: %s", initialMsg)

	// Send registration immediately
	registration := Message{
		Type: "register",
		ID:   clientID,
	}
	registrationJSON, _ := json.Marshal(registration)
	if err := websocket.Message.Send(ws, string(registrationJSON)); err != nil {
		log.Fatal("Registration error:", err)
	}

	// Set up message handler
	go func() {
		defer close(done)
		for {
			var rawMessage string
			if err := websocket.Message.Receive(ws, &rawMessage); err != nil {
				log.Println("Read error:", err)
				return
			}

			var message Message
			if err := json.Unmarshal([]byte(rawMessage), &message); err != nil {
				log.Printf("JSON parse error: %v", err)
				continue
			}

			log.Printf("Received: %s", rawMessage)

			// Handle capability commands
			if message.Type == "command" && message.Command != "" {
				response, err := handleCapabilityCommand(message.Command, message.Files)

				responseMsg := Message{
					Type:      "response",
					Message:   response,
					Command:   message.Command,
					Timestamp: time.Now().Format(time.RFC3339),
					ID:        clientID,
					TargetID:  message.ID,
				}

				if err != nil {
					responseMsg.Type = "error"
					responseMsg.Message = err.Error()
				}

				responseJSON, err := json.Marshal(responseMsg)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}

				if err := websocket.Message.Send(ws, string(responseJSON)); err != nil {
					log.Printf("Error sending response: %v", err)
				}
			}
		}
	}()

	// Handle shutdown
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, shutting down...")
			ws.Close()
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
