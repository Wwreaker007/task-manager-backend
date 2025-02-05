package contracts

import "task-manager-backend/data/common"

type GetAllUsersResponse struct{
	Status 		string				`json:"status"`
	Users		[]common.User		`json:"users"`
}