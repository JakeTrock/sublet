package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/lithdew/flatend"
)

func main() {
	// Create Flatend node
	node := &flatend.Node{
		Services: map[string]flatend.Handler{
			"init_configuration": handleInitConfiguration,
			"list_nix_files":     handleListNixFiles,
			"get_nix_contents":   handleGetNixFilesContents,
			"set_nix_contents":   handleSetNixFilesContents,
			"run_dry_build":      handleRunDryBuild,
			"run_test":           handleRunTest,
			"run_switch":         handleRunSwitch,
			"liveness_check":     handleLivenessCheck,
			"fetch_local_url":    handleFetchLocalUrl,
			"ssh":                handleSSHWebSocket,
		},
	}

	fmt.Println("Server started at :9000")
	if err := node.Start("127.0.0.1:9000"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	node.Shutdown()
}
