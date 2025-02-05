package common

type TASK_STATUS string

const (
	PENDING 	TASK_STATUS	= "PENDING"
	APPROVED 	TASK_STATUS = "APPROVED"
	REJECTED	TASK_STATUS = "REJECTED"
)

const (
	EMAIL_SUBJECT_APPROVER = "TASKID [%d] : ACTION REQUIERED"
	EMAIL_BODY_APPROVER = "Please review the task raised by userID : %d."

	EMAIL_SUBJECT_APPROVED = "TASKID [%d] : TASK APPROVED"
	EMAIL_BODY_APPROVED = "Task has been approved by [%d, %d, %d]. CreatorID [%d]"
)

type Task struct{
	TaskData
	ID 				uint64				`json:"id"`
	TaskStatus 		TASK_STATUS			`json:"task_status"`
	CreatorID		uint64				`json:"creator_id"`
	Comments		[]TaskComment		`json:"comments"`
	Approvals		[]TaskApproval		`json:"approvals"`
	CreatedAt		string				`json:"created_at"`
}

type TaskData struct{
	Title 			string			`json:"title"`
	Description 	string			`json:"description"`
	Approvers		[]uint64		`json:"approvers"`
}

type TaskApproval struct{
	TaskApprovalData
	ID 				uint64			`json:"id"`
	ApproverID		uint64			`json:"approver_id"`
	ApprovedAt		string			`json:"approved_at"`
}

type TaskApprovalData struct{
	TaskID			uint64			`json:"task_id"`
}

type TaskComment struct{
	TaskCommentData
	ID 				uint64			`json:"id"`
	Commentor		uint64			`json:"user_id"`
	CommentedAt		string			`json:"commented_at"`
}

type TaskCommentData struct{
	TaskID			uint64			`json:"task_id"`
	Comment			string			`json:"comment"`
}

type NotificationPayload struct{
	Recipient 		string
	Subject			string
	Body 			string
}