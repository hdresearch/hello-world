package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	db    *pgxpool.Pool
	token string
}

type countResponse struct {
	Count int64 `json:"count"`
}

func required(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("%s is required", name)
	}
	return value
}

func (a *application) authorized(request *http.Request) bool {
	actual := request.Header.Get("Authorization")
	expected := "Bearer " + a.token
	return len(actual) == len(expected) && subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) == 1
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}

func (a *application) health(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(request.Context(), time.Second)
	defer cancel()
	if err := a.db.Ping(ctx); err != nil {
		http.Error(writer, "unavailable", http.StatusServiceUnavailable)
		return
	}
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("ok\n"))
}

func (a *application) visits(writer http.ResponseWriter, request *http.Request) {
	if !a.authorized(request) {
		writeJSON(writer, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), 3*time.Second)
	defer cancel()
	var count int64
	var err error
	switch request.Method {
	case http.MethodGet:
		err = a.db.QueryRow(ctx, "SELECT count FROM visits WHERE singleton = TRUE").Scan(&count)
	case http.MethodPost:
		err = a.db.QueryRow(ctx, "UPDATE visits SET count = count + 1 WHERE singleton = TRUE RETURNING count").Scan(&count)
	default:
		writer.Header().Set("Allow", "GET, POST")
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		log.Printf("visit query failed: %v", err)
		writeJSON(writer, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}
	writeJSON(writer, http.StatusOK, countResponse{Count: count})
}

func main() {
	config, err := pgxpool.ParseConfig("postgres://postgres@127.0.0.1:5432/webstack?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	config.ConnConfig.Password = required("POSTGRES_PASSWORD")
	config.MaxConns = 4
	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{db: db, token: required("BACKEND_TOKEN")}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", app.health)
	mux.HandleFunc("/visits", app.visits)
	server := &http.Server{
		Addr:              ":80",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	log.Print("database gateway listening on :80")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
