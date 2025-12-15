package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/JpUnique/goqueue/internal/service"
)

type JobHandler struct {
	svc *service.JobService
}

func NewJobHandler(svc *service.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

func (h *JobHandler) Register(r chi.Router) {
	r.Post("/jobs", h.create)
	r.Get("/jobs/{id}", h.get)
}

type createJobReq struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (h *JobHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createJobReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	id, err := h.svc.CreateJob(r.Context(), req.Type, req.Payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *JobHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	job, err := h.svc.GetJob(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	json.NewEncoder(w).Encode(job)
}
