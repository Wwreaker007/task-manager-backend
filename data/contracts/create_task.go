package contracts

import "task-manager-backend/data/common"

type CreateTaskRequest struct{
	common.TaskData
}

type CreateTaskResponse struct{
	Status		string		`json:"status"`
	TaskID		uint64		`json:"task_id"`
}