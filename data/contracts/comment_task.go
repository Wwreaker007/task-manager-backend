package contracts

import "task-manager-backend/data/common"

type CommentTaskRequest struct{
	common.TaskCommentData
}

type CommentTaskResponse struct {
	Status			string		`json:"status"`
}
