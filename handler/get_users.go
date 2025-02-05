package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"task-manager-backend/types"
	"time"
)

type GetUsersHandler struct{
	usm 		types.UserService
}

func NewGetUsersHandler(service types.UserService) *GetUsersHandler{
	return &GetUsersHandler{
		usm: service,
	}
}

func (h *GetUsersHandler) Handle(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	response := h.usm.GetAllUsers(ctx)
	json.NewEncoder(w).Encode(response)
}