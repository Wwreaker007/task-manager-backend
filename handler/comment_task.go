package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"task-manager-backend/data/contracts"
	"task-manager-backend/types"
	"time"
)

type CommentTaskHandler struct{
	tsm 		types.TaskService
}

func NewCommentTaskHandler(service types.TaskService) *CommentTaskHandler{
	return &CommentTaskHandler{
		tsm: service,
	}
}

func (h *CommentTaskHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	var req contracts.CommentTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.tsm.CommentTask(ctx, req)
	json.NewEncoder(w).Encode(response)
}