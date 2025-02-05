package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"task-manager-backend/data/contracts"
	"task-manager-backend/types"
	"time"
)

type LoginHandler struct{
	usm 		types.UserService
}

func NewLoginHandler(service types.UserService) *LoginHandler{
	return &LoginHandler{
		usm: service,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	var req contracts.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.usm.Login(ctx, req)
	json.NewEncoder(w).Encode(response)
}