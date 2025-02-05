package psql

import (
	"context"
	"database/sql"
	"log"
	"task-manager-backend/data/common"
)

const (
	CREATE_APPROVAL_QUERY = "INSERT INTO task_approvals (task_id, approver_id) VALUES ($1, $2)"
	GET_APPROVALS_BY_TASKID_QUERY = "SELECT * FROM task_approvals WHERE task_id = $1"
	GET_APPROVALS_BY_USERID_QUERY = "SELECT * FROM task_approvals WHERE approver_id = $1"
)

type ApprovalsRepositoryManager struct{
	db 		*sql.DB
}

func NewApprovalsRepositoryManager(client *sql.DB) *ApprovalsRepositoryManager{
	return &ApprovalsRepositoryManager{
		db: client,
	}
}

func (arm *ApprovalsRepositoryManager) CreateApproval(ctx context.Context, data common.TaskApprovalData) error {
	userID := ctx.Value(common.USER_ID_KEY)
	_, err := arm.db.Exec(CREATE_APPROVAL_QUERY, data.TaskID, userID)
	if err != nil {
		log.Println("Error in creating approval : ", err)
		return err
	}
	return nil
}

func (arm *ApprovalsRepositoryManager) GetApprovalsByTaskID(ctx context.Context, taskID uint64) ([]common.TaskApproval, error) {
	rows, err := arm.db.Query(GET_APPROVALS_BY_TASKID_QUERY, taskID)
	if err != nil {
		log.Println("Error in fetching approvals : ", err)
		return nil, err
	}
	return scanApprovalRecords(rows)
}

func (arm *ApprovalsRepositoryManager) GetApprovalsByUserID(ctx context.Context) ([]common.TaskApproval, error) {
	userID := ctx.Value(common.USER_ID_KEY)
	rows, err := arm.db.Query(GET_APPROVALS_BY_TASKID_QUERY, userID)
	if err != nil {
		log.Println("Error in fetching approvals : ", err)
		return nil, err
	}
	return scanApprovalRecords(rows)
}

func scanApprovalRecords(rows *sql.Rows) ([]common.TaskApproval, error){
	var taskApprovals []common.TaskApproval
	for rows.Next() {
		var taskApproval common.TaskApproval
		err := rows.Scan(&taskApproval.ID, &taskApproval.TaskID, &taskApproval.ApproverID, &taskApproval.ApprovedAt)
		if err != nil {
			log.Println("Error in scanning approvals : ", err)
			return nil, err
		}
		taskApprovals = append(taskApprovals, taskApproval)
	}
	return taskApprovals, nil
}

