package http

import (
	"encoding/json"
	"net/http"

	"repo-stat/api/internal/usecase"
)

type Handler struct {
	mux *http.ServeMux
}

func NewHandler(pingUC *usecase.Ping, repoUC *usecase.GetRepository) *Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/ping", pingHandler(pingUC))
	mux.HandleFunc("GET /repositories/{owner}/{repo}", repositoryHandler(repoUC))

	return &Handler{mux: mux}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		return
	}
}

func pingHandler(pingUC *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		result, statusCode := pingUC.Execute(r.Context())

		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"status":"error","message":"failed to encode response"}`))
			if err != nil {
				return
			}
			return
		}
	}
}

func repositoryHandler(repoUC *usecase.GetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		owner := r.PathValue("owner")
		repo := r.PathValue("repo")
		url := "https://github.com/" + owner + "/" + repo

		repository, err := repoUC.Execute(r.Context(), url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": err.Error(),
			})
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(repository)
		if err != nil {
			return
		}
	}
}
