package main

import (
	"fmt"
	"net/http"
)

type config struct {
	Host         string `toml:"host"`
	Port         uint   `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	IdentityFile string `toml:"identity_file"`
}

func main() {
	http.HandleFunc("/initConfiguration", handleInitConfiguration)
	http.HandleFunc("/listNixFiles", handleListNixFiles)
	http.HandleFunc("/getNixFilesContents", handleGetNixFilesContents)
	http.HandleFunc("/setNixFilesContents", handleSetNixFilesContents)
	http.HandleFunc("/runDryBuild", handleRunDryBuild)
	http.HandleFunc("/runTest", handleRunTest)
	http.HandleFunc("/runSwitch", handleRunSwitch)
	http.HandleFunc("/livenessCheck", handleLivenessCheck)
	http.HandleFunc("/fetchLocalUrl", handleFetchLocalUrl)

	handler := &sshHandler{
		addr:   "127.0.0.1",
		user:   "root",            //TODO: do not use root
		secret: "8s9f8ds9f9d8fds", //TODO: randomize, do not hardcode
	}
	http.HandleFunc("/ssh", handler.webSocket)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
