package contracts

import "task-manager-backend/data/common"

type ApproveTaskRequest struct {
	common.TaskApprovalData
}

type ApproveTaskResponse struct{
	common.Task
	Status		string		`json:"status"`
}