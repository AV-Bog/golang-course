package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, statusCode := ping.Execute(r.Context())

		response := dto.PingResponseDTO{
			Status:   result.Status,
			Services: make([]dto.ServiceStatusDTO, len(result.Services)),
		}

		for i, svc := range result.Services {
			response.Services[i] = dto.ServiceStatusDTO{
				Name:   svc.Name,
				Status: svc.Status,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
