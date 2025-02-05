package contracts

import "task-manager-backend/data/common"

type SignUpRequest struct{
	common.UserData
}

type SignUpResponse struct{
	Status		string		`json:"status"`
}