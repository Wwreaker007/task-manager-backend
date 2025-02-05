package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"task-manager-backend/data/contracts"
	"task-manager-backend/types"
	"time"
)

type ApproveTaskHandler struct{
	tsm 		types.TaskService
}

func NewApproveTaskHandler(service types.TaskService) *ApproveTaskHandler{
	return &ApproveTaskHandler{
		tsm: service,
	}
}

func (h *ApproveTaskHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	var req contracts.ApproveTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.tsm.ApproveTask(ctx, req)
	json.NewEncoder(w).Encode(response)
}

