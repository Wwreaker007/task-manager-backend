package types

import (
	"context"
	"task-manager-backend/data/common"
	"task-manager-backend/data/contracts"
)

type ApprovalsRepository interface{
	CreateApproval(context.Context, common.TaskApprovalData) error
	GetApprovalsByTaskID(context.Context, uint64) ([]common.TaskApproval, error) 
	GetApprovalsByUserID(context.Context) ([]common.TaskApproval, error)
}

type CommentsRepository interface{
	CreateComment(context.Context, common.TaskCommentData) error
	GetCommentsByTaskID(context.Context, uint64) ([]common.TaskComment, error)
	GetCommentsByUserID(context.Context) ([]common.TaskComment, error)
}

type TasksRepository interface{
	CreateTask(context.Context, common.TaskData) (uint64, error)
	GetTaskByTaskID(context.Context, uint64) (common.Task, error)
	GetTasksByUserID(context.Context) ([]common.Task, error)
	UpdateTaskByTaskID(context.Context, common.Task) (common.Task, error)
}

type UsersRepository interface{
	CreateUser(context.Context, common.UserData) error
	GetUserByUserID(context.Context, uint64) (common.User, error)
	GetUserByName(context.Context, string) (common.User, error)
	GetAllUsers(context.Context) ([]common.User, error)
}

type TaskService interface{
	CreateTask(context.Context, contracts.CreateTaskRequest) contracts.CreateTaskResponse
	ApproveTask(context.Context, contracts.ApproveTaskRequest) contracts.ApproveTaskResponse
	CommentTask(context.Context, contracts.CommentTaskRequest) contracts.CommentTaskResponse
	GetTask(context.Context, uint64) contracts.GetTaskResponse 
}

type UserService interface{
	SignUp(context.Context, contracts.SignUpRequest) contracts.SignUpResponse
	Login(context.Context, contracts.LoginRequest) contracts.LoginResponse
	GetAllUsers(context.Context) contracts.GetAllUsersResponse
}

type NotificationService interface{
	SendNotification(context.Context, []common.NotificationPayload)
}