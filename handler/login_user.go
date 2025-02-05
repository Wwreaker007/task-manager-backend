package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"task-manager-backend/data/contracts"
	"task-manager-backend/types"
	"time"
)

type SignUpHandler struct{
	usm 		types.UserService
}

func NewSignUpHandler(service types.UserService) *SignUpHandler{
	return &SignUpHandler{
		usm: service,
	}
}

func (h *SignUpHandler) Handle(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	var req contracts.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.usm.SignUp(ctx, req)
	json.NewEncoder(w).Encode(response)
}