package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"task-manager-backend/data/contracts"
	"task-manager-backend/types"
	"time"
)

type CreateTaskHandler struct{
	tsm 		types.TaskService
}

func NewCreateTaskHandler(service types.TaskService) *CreateTaskHandler{
	return &CreateTaskHandler{
		tsm: service,
	}
}

func (h *CreateTaskHandler) Handle(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	var req contracts.CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.tsm.CreateTask(ctx, req)
	json.NewEncoder(w).Encode(response)
}