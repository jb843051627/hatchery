package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jb843051627/hatchery/internal/cache"
	"github.com/jb843051627/hatchery/internal/handler"
	"github.com/jb843051627/hatchery/internal/service"
	"github.com/jb843051627/hatchery/internal/store"
)

func main() {
	dbPath := os.Getenv("HATCHERY_DB")
	if dbPath == "" {
		dbPath = filepath.Join(os.TempDir(), "hatchery.db")
	}
	st, err := store.NewStore(dbPath)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}
	defer st.Close()
	rc := cache.NewReadingCache()
	svc := service.NewService(st, rc)
	h := handler.NewHandler(svc)
	mux := h.Routes()
	addr := os.Getenv("HATCHERY_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("hatchery listening on %s (db=%s)", addr, dbPath)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
