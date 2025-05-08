package api

import (
	"healthmonitor/internal/db"
	"net/http"

	"encoding/json"

	"github.com/rs/zerolog/log"
)

type APIHandler struct {
	dbh *db.DB
}

func New(dbh *db.DB) (*APIHandler, error) {
	return &APIHandler{
		dbh: dbh,
	}, nil
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /counts", h.handleGetCounts)
	mux.HandleFunc("GET /devices", h.handleGetDevices)
	mux.ServeHTTP(w, r)
}

func (h *APIHandler) handleGetCounts(w http.ResponseWriter, r *http.Request) {
	counts, err := h.dbh.GetDeviceCounts()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get device counts")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(counts)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *APIHandler) handleGetDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.dbh.GetDevices()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get devices")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(devices)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
