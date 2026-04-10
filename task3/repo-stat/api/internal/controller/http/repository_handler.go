package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"repo-stat/api/internal/usecase"
)

func NewRepositoryHandler(log *slog.Logger, getRepo *usecase.GetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "url parameter is required"})
			return
		}

		resp, err := getRepo.Execute(r.Context(), url)
		if err != nil {
			log.Error("failed to get repository", "error", err)
			w.Header().Set("Content-Type", "application/json")

			errMsg := err.Error()
			statusCode := http.StatusInternalServerError

			if strings.Contains(errMsg, "not a GitHub URL") ||
				strings.Contains(errMsg, "InvalidArgument") ||
				strings.Contains(errMsg, "invalid repository URL") ||
				strings.Contains(errMsg, "invalid format") {
				statusCode = http.StatusBadRequest
			}

			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp.Repository); err != nil {
			log.Error("failed to encode response", "error", err)
		}
	}
}
