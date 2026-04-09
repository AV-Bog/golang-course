package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/usecase"
)

func NewRepositoryHandler(log *slog.Logger, getRepo *usecase.GetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "url parameter is required", http.StatusBadRequest)
			return
		}

		resp, err := getRepo.Execute(r.Context(), url)
		if err != nil {
			log.Error("failed to get repository", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", "error", err)
		}
	}
}
