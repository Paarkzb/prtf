package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type StreamKey struct {
	Key string `json:"key"`
}

var (
	activeStreams = make(map[string]bool)
	mu            sync.Mutex
)

type Recording struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Date     string `json:"date"`
	Duration string `json:"duration"`
}

func listRecordingsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("/var/vod")
	if err != nil {
		http.Error(w, "failed to read recordings", http.StatusInternalServerError)
		return
	}

	var recordings []Recording
	for _, f := range files {
		fileInfo, err := f.Info()
		if err != nil {
			http.Error(w, "failed to read recordings info", http.StatusInternalServerError)
			return
		}

		if filepath.Ext(f.Name()) == ".mp4" {
			recordings = append(recordings, Recording{
				Name: f.Name(),
				Path: f.Name(),
				Date: fileInfo.ModTime().Format("2006-01-02 15:14"),
			})
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(recordings)
}

func authStreamHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("name")
	log.Println(r.URL)
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	// Проверка валидности ключа (здесь пример, замените на свою логику)
	if isValidStreamKey(key) {
		mu.Lock()
		activeStreams[key] = true
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Invalid stream key", http.StatusUnauthorized)
}

func isValidStreamKey(key string) bool {
	// Пример: проверка ключа в базе данных
	return key == "test123"
}

func listStreamsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(activeStreams)
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("pong")
}

func main() {
	http.HandleFunc("/api/auth/stream", authStreamHandler)
	http.HandleFunc("/api/streams", listStreamsHandler)
	http.HandleFunc("/", pong)

	http.HandleFunc("/api/recordings", listRecordingsHandler)
	http.Handle("/vod/", http.StripPrefix("/vod/", http.FileServer(http.Dir("/var/vod"))))

	fmt.Println("Server running on :8092")
	log.Fatal(http.ListenAndServe(":8092", nil))
}
