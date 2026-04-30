package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"repo-stat/api/internal/usecase"
)

type RepositoryHandler struct {
	getRepoUC *usecase.GetRepository
}

func NewRepositoryHandler(getRepoUC *usecase.GetRepository) *RepositoryHandler {
	return &RepositoryHandler{
		getRepoUC: getRepoUC,
	}
}

func (h *RepositoryHandler) GetRepository(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/repositories/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid repository path", http.StatusBadRequest)
		return
	}

	url := "https://github.com/" + parts[0] + "/" + parts[1]

	repo, err := h.getRepoUC.Execute(r.Context(), url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"full_name":   repo.FullName,
		"description": repo.Description,
		"stars":       repo.Stars,
		"forks":       repo.Forks,
		"created_at":  repo.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"language":    repo.Language,
	})
	if err != nil {
		return
	}
}
