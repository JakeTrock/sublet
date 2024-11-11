package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/initConfiguration", handleInitConfiguration)
	http.HandleFunc("/listNixFiles", handleListNixFiles)
	http.HandleFunc("/getNixFilesContents", handleGetNixFilesContents)
	http.HandleFunc("/setNixFilesContents", handleSetNixFilesContents)
	http.HandleFunc("/runDryBuild", handleRunDryBuild)
	http.HandleFunc("/runTest", handleRunTest)
	http.HandleFunc("/runSwitch", handleRunSwitch)
	http.HandleFunc("/livenessCheck", handleLivenessCheck)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func handleInitConfiguration(w http.ResponseWriter, r *http.Request) {
	err := initConfiguration()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleListNixFiles(w http.ResponseWriter, r *http.Request) {
	files, err := listNixFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(files)
}

func handleGetNixFilesContents(w http.ResponseWriter, r *http.Request) {
	contents, err := getNixFilesContents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(contents))
}

func handleSetNixFilesContents(w http.ResponseWriter, r *http.Request) {
	var files map[string]string
	err := json.NewDecoder(r.Body).Decode(&files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = setNixFilesContents(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleRunDryBuild(w http.ResponseWriter, r *http.Request) {
	output, err := runDryBuild()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(output))
}

func handleRunTest(w http.ResponseWriter, r *http.Request) {
	output, err := runTest()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(output))
}

func handleRunSwitch(w http.ResponseWriter, r *http.Request) {
	output, err := runSwitch()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(output))
}

func handleLivenessCheck(w http.ResponseWriter, r *http.Request) {
	output, err := livenessCheck()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(output))
}
