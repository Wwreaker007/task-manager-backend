package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"task-manager-backend/types"
	"time"

	"github.com/gorilla/mux"
)

const (
	task_ID = "id"
)

type GetTaskHandler struct{
	tsm 		types.TaskService
}

func NewGetTaskHandler(service types.TaskService) *GetTaskHandler{
	return &GetTaskHandler{
		tsm: service,
	}
}

func (h *GetTaskHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	// Extract the path parametet
	vars := mux.Vars(r)
	taskID := vars[task_ID] 

	// Convert integer string to string
	intTaskID, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.tsm.GetTask(ctx, uint64(intTaskID))
	json.NewEncoder(w).Encode(response)
}