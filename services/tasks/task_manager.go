package tasks

import (
	"context"
	"fmt"
	"task-manager-backend/data/common"
	"task-manager-backend/data/contracts"
	"task-manager-backend/types"
)

type TaskServiceManager struct{
	tasksRepo	 		types.TasksRepository
	commentsRepo		types.CommentsRepository
	approvalRepo		types.ApprovalsRepository
	userRepo			types.UsersRepository
	notify 				types.NotificationService
}

func NewTaskServiceManager(taskDb types.TasksRepository, commentDb types.CommentsRepository, approvalDb types.ApprovalsRepository, userDb types.UsersRepository, notify types.NotificationService) *TaskServiceManager {
	return &TaskServiceManager{
		tasksRepo: taskDb,
		commentsRepo: commentDb,
		approvalRepo: approvalDb,
		userRepo: userDb,
		notify: notify,
	}
}

func (tsm *TaskServiceManager) CreateTask(ctx context.Context, r contracts.CreateTaskRequest) contracts.CreateTaskResponse {
	taskID, err := tsm.tasksRepo.CreateTask(ctx, r.TaskData)
	if err != nil {
		return contracts.CreateTaskResponse{
			Status: "FAILED : " + err.Error(),
		}
	}

	// Create a goroutine to handle the notification handling
	go func() {
		// Build notification payload for the approvers
		var messages []common.NotificationPayload
		for _, user := range r.Approvers {
			userData, _ := tsm.userRepo.GetUserByUserID(ctx, user)
			msg := common.NotificationPayload{
				Recipient: userData.Email,
				Subject: fmt.Sprintf(common.EMAIL_SUBJECT_APPROVER, taskID),
				Body: fmt.Sprintf(common.EMAIL_BODY_APPROVER, ctx.Value(common.USER_ID_KEY)),
			}
			messages = append(messages, msg)
		}

		// Notify the approvers
		tsm.notify.SendNotification(ctx, messages)
	}()

	return contracts.CreateTaskResponse{
		Status: "SUCCESS",
		TaskID: taskID,
	}
}

func (tsm *TaskServiceManager) ApproveTask(ctx context.Context, r contracts.ApproveTaskRequest) contracts.ApproveTaskResponse {
	// Check if the task needs an approval
	task, err := tsm.tasksRepo.GetTaskByTaskID(ctx, r.TaskID)
	if err != nil {
		return contracts.ApproveTaskResponse{
			Status: "FAILED : " + err.Error(),
		}
	}

	// Append the comments to the final response
	comments, err := tsm.commentsRepo.GetCommentsByTaskID(ctx, r.TaskID)
	if err != nil {
		return contracts.ApproveTaskResponse{
			Status: "P-SUCCESS : Error in fetching task comments : " +  err.Error(),
			Task: task,
		}
	}

	// Add the fetched comments in the task instance
	task.Comments = comments

	// Check if the task is already approved
	if task.TaskStatus == common.APPROVED {
		return contracts.ApproveTaskResponse{
			Status: "P-SUCCESS : Task is already approved",
			Task: task,
		}
	}

	// Creator cannot approve its own task
	if task.CreatorID == ctx.Value(common.USER_ID_KEY) {
		return contracts.ApproveTaskResponse{
			Status: "FAILED : Cannot approve your own task !",
			Task: task,
		}
	}

	// Proceed with recording this new approval after passing all the checks
	err = tsm.approvalRepo.CreateApproval(ctx, r.TaskApprovalData)
	if err != nil {
		return contracts.ApproveTaskResponse{
			Status: "FAILED : " + err.Error(),
			Task: task,
		}
	}

	// Retrive all the approvals for the taskID
	approvals, err := tsm.approvalRepo.GetApprovalsByTaskID(ctx, r.TaskID)
	if err != nil {
		return contracts.ApproveTaskResponse{
			Status: "FAILED : " + err.Error(),
			Task: task,
		}
	}

	// Assign the approvals 
	task.Approvals = approvals

	// Count the total number of approvals, if approvals == 3, approve the task
	if len(approvals) < 3 {
		return contracts.ApproveTaskResponse{
			Status: fmt.Sprintf("SUCCESS : Waiting for %d more approvals", 3 - len(approvals)),
			Task: task,
		}
	}

	// The task has been approved by 3 approvers and update in taskDB
	task.TaskStatus = common.APPROVED
	_, err = tsm.tasksRepo.UpdateTaskByTaskID(ctx, task)
	if err != nil {
		return contracts.ApproveTaskResponse{
			Status: "Error in updating task status :" + err.Error(),
		}
	}

	// Create a goroutine to handle the notification handling
	go func() {
		// Create the memberList involved in the task
		var membersList []uint64
		for _, approval := range approvals {
			membersList = append(membersList, approval.ApproverID)
		}
		membersList = append(membersList, task.CreatorID)

		// Build notification payload for the approvers
		var messages []common.NotificationPayload
		for _, user := range membersList {
			userData, _ := tsm.userRepo.GetUserByUserID(ctx, user)
			msg := common.NotificationPayload{
				Recipient: userData.Email,
				Subject: fmt.Sprintf(common.EMAIL_SUBJECT_APPROVED, r.TaskID),
				Body: fmt.Sprintf(common.EMAIL_BODY_APPROVED, membersList[0], membersList[1], membersList[2], membersList[3]),
			}
			messages = append(messages, msg)
		}

		// Notify the members
		tsm.notify.SendNotification(ctx, messages)
	}()

	// Return the complete task response
	return contracts.ApproveTaskResponse{
		Status: "SUCCESS : Task approved and notified",
		Task: task,
	}
}

func (tsm *TaskServiceManager) CommentTask(ctx context.Context, r contracts.CommentTaskRequest) contracts.CommentTaskResponse {
	err := tsm.commentsRepo.CreateComment(ctx, r.TaskCommentData)
	if err != nil {
		return contracts.CommentTaskResponse{
			Status: "FAILED : " + err.Error(),
		}
	}
	return contracts.CommentTaskResponse{
		Status: "SUCCESS",
	}
}

func (tsm *TaskServiceManager) GetTask(ctx context.Context, taskID uint64) contracts.GetTaskResponse {
	// Get the task using the taskID
	task, err := tsm.tasksRepo.GetTaskByTaskID(ctx, taskID)
	if err != nil {
		return contracts.GetTaskResponse{
			Status: "FAILED : " + err.Error(),
		}
	}

	// Fetch the comments on the basis on taskID
	comments, err := tsm.commentsRepo.GetCommentsByTaskID(ctx, taskID)
	if err != nil {
		return contracts.GetTaskResponse{
			Status: "P-SUCCESS : Error in fetching task comments : " +  err.Error(),
			Task: task,
		}
	}

	// Add the fetched comments in the task instance
	task.Comments = comments

	// Fetch the approvals on the basis on taskID
	approvals, err := tsm.approvalRepo.GetApprovalsByTaskID(ctx, taskID)
	if err != nil {
		return contracts.GetTaskResponse{
			Status: "P-SUCCESS : Error in fetching task approvals : " +  err.Error(),
			Task: task,
		}
	}

	// Add the fetched approval sin the task instance
	task.Approvals = approvals

	// Return the collated instance
	return contracts.GetTaskResponse{
		Status: "SUCCESS",
		Task: task,
	}
}
