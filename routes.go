package main

import (
	"encoding/json"
	"net/http"
)

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
