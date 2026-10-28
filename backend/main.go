package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/YuminosukeSato/pyproc/pkg/pyproc"
)

var pool *pyproc.Pool

func creat_pool() {
	var err error

	pool, err = pyproc.NewPool(pyproc.PoolOptions{
		Config: pyproc.PoolConfig{
			Workers:     8,
			MaxInFlight: 10,
		},
		WorkerConfig: pyproc.WorkerConfig{
			SocketPath:   "/tmp/worker.sock",
			PythonExec:   "python3",
			WorkerScript: "./python-worker/worker.py",
			StartTimeout: 30 * time.Second,
		},
	}, nil)
	
	if err != nil {
		log.Fatal("failed to create pool: ", err)
	}

	if err := pool.Start(context.Background()); err != nil {
		log.Fatal("failed to start pool: ", err)
	}

	log.Println("Pool was created and started successfully")
}

type FromPy struct {
	TargetAI    string `json:"target_ai"`
	BaseInput   string `json:"instruction"` 
	Strategy    string `json:"strategy"`
	AuthToken   string `json:"auth_token"`
	ModleName   string `json:"modle_name"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		return
	}

	var req FromPy
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode: ", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input := map[string]interface{}{
		"target_ai":   req.TargetAI,
		"instruction": req.BaseInput,
		"strategy":    req.Strategy,
		"auth_token":  req.AuthToken,
		"modle_name":  req.ModleName,
	}

	var output map[string]interface{}

	if err := pool.Call(context.Background(), "combine", input, &output); err != nil {
		log.Println("failed to call pool: ", err)
		http.Error(w, "Python call failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}


    func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    initPool()
    defer pool.Shutdown(context.Background())
    http.HandleFunc("/scan", handle)
    log.Println("Go server running on 0.0.0.0:" + port)
    log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}