package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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

	cached, err := rdb.Get(context.Background(), "streams").Bytes()
	if err == nil {
		w.Write(cached)
		return
	}

	mu.Lock()
	rdb.Set(context.Background(), "streams", activeStreams, 10*time.Second)

	defer mu.Unlock()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(activeStreams)
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("pong")
}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "prtf-stream-redis:6379",
	Password: "redis",
	DB:       0,
})

func recordMetrics() {
	go func() {
		for {
			promActiveStreams.Add(float64(len(activeStreams)))

			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	promStreamRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "stream_requests_total",
			Help: "Total number of stream requests",
		},
	)
	promActiveStreams = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_streams",
			Help: "Current active streams",
		},
	)
)

func init() {
	prometheus.MustRegister(promStreamRequests, promActiveStreams)
}

func main() {
	recordMetrics()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/api/auth/stream", authStreamHandler)
	http.HandleFunc("/api/streams", listStreamsHandler)
	http.HandleFunc("/", pong)

	http.HandleFunc("/api/recordings", listRecordingsHandler)
	http.Handle("/vod/", http.StripPrefix("/vod/", http.FileServer(http.Dir("/var/vod"))))

	sugar.Infow("Server running on :8092")
	sugar.Fatal(http.ListenAndServe(":8092", nil))
}
