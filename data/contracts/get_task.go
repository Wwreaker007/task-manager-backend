package contracts

import "task-manager-backend/data/common"

type GetTaskResponse struct{
	common.Task
	Status 			string 			`json:"status"`
}