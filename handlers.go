package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lithdew/flatend"
	"golang.org/x/crypto/ssh"
)

// Add these constants and types at the top of handlers.go
const (
	messageWait    = 10 * time.Second
	maxMessageSize = 512
)

var terminalModes = ssh.TerminalModes{
	ssh.ECHO:          1,     // enable echoing
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}

type windowSize struct {
	High  int `json:"high"`
	Width int `json:"width"`
}

type sshSession struct {
	ctx      *flatend.Context
	client   *ssh.Client
	sess     *ssh.Session
	sessIn   io.WriteCloser
	sessOut  io.Reader
	closeSig chan struct{}
}

func (s *sshSession) getWindowSize() (*windowSize, error) {
	body, err := io.ReadAll(s.ctx.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read window size: %w", err)
	}

	var wdSize windowSize
	if err := json.Unmarshal(body, &wdSize); err != nil {
		return nil, fmt.Errorf("invalid window size format: %w", err)
	}
	return &wdSize, nil
}

func (s *sshSession) handleOutput() {
	defer func() { s.closeSig <- struct{}{} }()

	buf := make([]byte, maxMessageSize)
	for {
		time.Sleep(10 * time.Millisecond)
		n, err := s.sessOut.Read(buf)
		if n > 0 {
			if _, err := s.ctx.Write(buf[:n]); err != nil {
				log.Printf("Error writing to client: %v", err)
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from SSH: %v", err)
			}
			return
		}
	}
}

func (s *sshSession) handleInput() {
	defer func() { s.closeSig <- struct{}{} }()

	buf := make([]byte, maxMessageSize)
	for {
		n, err := s.ctx.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from client: %v", err)
			}
			return
		}

		// Check if this is a window resize message
		var wdSize windowSize
		if err := json.Unmarshal(buf[:n], &wdSize); err == nil {
			if err := s.sess.WindowChange(wdSize.High, wdSize.Width); err != nil {
				log.Printf("Error changing window size: %v", err)
			}
			continue
		}

		// Otherwise treat as regular input
		if _, err := s.sessIn.Write(buf[:n]); err != nil {
			log.Printf("Error writing to SSH: %v", err)
			return
		}
	}
}

// Helper function to write JSON response
func writeJSON(ctx *flatend.Context, data interface{}) {
	response, _ := json.Marshal(data)
	ctx.WriteHeader("Content-Type", "application/json")
	ctx.Write(response)
}

// Helper function to write error response
func writeError(ctx *flatend.Context, err error) {
	writeJSON(ctx, map[string]string{
		"error": err.Error(),
	})
}

// Helper function to check if URL is local
func isLocalURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return strings.HasPrefix(u.Host, "localhost") || strings.HasPrefix(u.Host, "127.0.0.1")
}

func handleInitConfiguration(ctx *flatend.Context) {
	if err := initConfiguration(); err != nil {
		writeError(ctx, err)
		return
	}
	writeJSON(ctx, map[string]string{"status": "ok"})
}

func handleListNixFiles(ctx *flatend.Context) {
	files, err := listNixFiles()
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeJSON(ctx, map[string]interface{}{"files": files})
}

func handleGetNixFilesContents(ctx *flatend.Context) {
	contents, err := getNixFilesContents()
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeJSON(ctx, map[string]interface{}{"contents": contents})
}

func handleSetNixFilesContents(ctx *flatend.Context) {
	var req struct {
		Files map[string]string `json:"files"`
	}

	body, err := io.ReadAll(ctx.Body)
	if err != nil {
		writeError(ctx, fmt.Errorf("failed to read request body: %v", err))
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeError(ctx, fmt.Errorf("invalid request: %v", err))
		return
	}

	if err := setNixFilesContents(req.Files); err != nil {
		writeError(ctx, err)
		return
	}

	writeJSON(ctx, map[string]string{"status": "ok"})
}

func handleRunDryBuild(ctx *flatend.Context) {
	output, err := runDryBuild()
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeJSON(ctx, map[string]interface{}{"output": output})
}

func handleRunTest(ctx *flatend.Context) {
	output, err := runTest()
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeJSON(ctx, map[string]interface{}{"output": output})
}

func handleRunSwitch(ctx *flatend.Context) {
	output, err := runSwitch()
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeJSON(ctx, map[string]interface{}{"output": output})
}

func handleLivenessCheck(ctx *flatend.Context) {
	output, err := livenessCheck()
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeJSON(ctx, map[string]interface{}{"output": output})
}

func textLiveness(ctx *flatend.Context) {
	ctx.Write([]byte("Server is running"))
}

func handleFetchLocalUrl(ctx *flatend.Context) {
	var req struct {
		URL string `json:"url"`
	}

	body, err := io.ReadAll(ctx.Body)
	if err != nil {
		writeError(ctx, fmt.Errorf("failed to read request body: %v", err))
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeError(ctx, fmt.Errorf("invalid request body: %v", err))
		return
	}

	if !isLocalURL(req.URL) {
		writeError(ctx, fmt.Errorf("only localhost URLs are allowed"))
		return
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		writeError(ctx, fmt.Errorf("failed to fetch URL: %v", err))
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			ctx.WriteHeader(key, value)
		}
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		writeError(ctx, fmt.Errorf("error reading response: %v", err))
		return
	}

	ctx.Write(body)
}

func handleSSHWebSocket(ctx *flatend.Context) {
	// Extract connection parameters from headers
	addr := ctx.Headers["ssh-addr"]
	if addr == "" {
		addr = "127.0.0.1:22"
	}

	user := ctx.Headers["ssh-user"]
	if user == "" {
		user = "root" // TODO: do not use root by default
	}

	secret := ctx.Headers["ssh-secret"]
	keyfile := ctx.Headers["ssh-keyfile"]

	// Create SSH session
	session := &sshSession{
		ctx:      ctx,
		closeSig: make(chan struct{}, 1),
	}

	// Get initial window size
	wdSize, err := session.getWindowSize()
	if err != nil {
		writeError(ctx, fmt.Errorf("failed to get window size: %v", err))
		return
	}

	// Setup authentication
	var auth ssh.AuthMethod
	if secret != "" {
		auth = ssh.Password(secret)
	} else if keyfile != "" {
		key, err := os.ReadFile(keyfile)
		if err != nil {
			writeError(ctx, fmt.Errorf("failed to read key file: %v", err))
			return
		}
		privateKey, err := ssh.ParsePrivateKey(key)
		if err != nil {
			writeError(ctx, fmt.Errorf("failed to parse private key: %v", err))
			return
		}
		auth = ssh.PublicKeys(privateKey)
	} else {
		writeError(ctx, fmt.Errorf("no authentication method provided"))
		return
	}

	// Create SSH client
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: use proper host key verification
	}

	session.client, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		writeError(ctx, fmt.Errorf("failed to connect to SSH server: %v", err))
		return
	}
	defer session.client.Close()

	// Create new SSH session
	session.sess, err = session.client.NewSession()
	if err != nil {
		writeError(ctx, fmt.Errorf("failed to create SSH session: %v", err))
		return
	}
	defer session.sess.Close()

	// Setup pipes
	session.sessOut, err = session.sess.StdoutPipe()
	if err != nil {
		writeError(ctx, fmt.Errorf("failed to create stdout pipe: %v", err))
		return
	}

	session.sessIn, err = session.sess.StdinPipe()
	if err != nil {
		writeError(ctx, fmt.Errorf("failed to create stdin pipe: %v", err))
		return
	}
	defer session.sessIn.Close()

	// Request PTY
	if err := session.sess.RequestPty("xterm", wdSize.High, wdSize.Width, terminalModes); err != nil {
		writeError(ctx, fmt.Errorf("failed to request PTY: %v", err))
		return
	}

	// Start shell
	if err := session.sess.Shell(); err != nil {
		writeError(ctx, fmt.Errorf("failed to start shell: %v", err))
		return
	}

	log.Println("Started SSH session")
	defer log.Println("Closed SSH session")

	// Start input/output handlers
	go session.handleInput()
	go session.handleOutput()

	// Wait for session to end
	<-session.closeSig
}
