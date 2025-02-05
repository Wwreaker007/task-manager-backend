package psql

import (
	"context"
	"database/sql"
	"log"
	"task-manager-backend/data/common"

	"github.com/lib/pq"
)

const (
	CREATE_TASK_QUERY = "INSERT INTO tasks (title, description, task_status, creator_id, approvers) VALUES ($1, $2, $3, $4, $5)"
	RETURN_TASK_ID_QUERY = "INSERT INTO tasks (title, description, task_status, creator_id, approvers) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	GET_TASK_BY_TASKID_QUERY = "SELECT * FROM tasks WHERE id = $1"
	GET_TASKS_BY_USERID_QUERY = "SELECT * FROM tasks WHERE creator_id = $1"
	UPDATE_TASK_BY_TASKID = "UPDATE tasks SET title = $1, description = $2, task_status = $3, creator_id = $4, approvers = $5 WHERE id = $6"
)

type TasksRepositoryManager struct{
	db		*sql.DB
}

func NewTasksRepositoryManager(client *sql.DB) *TasksRepositoryManager{
	return &TasksRepositoryManager{
		db: client,
	}
}

func (trm *TasksRepositoryManager) CreateTask(ctx context.Context, data common.TaskData) (uint64, error){
	userID := ctx.Value(common.USER_ID_KEY)
	_, err := trm.db.Exec(CREATE_TASK_QUERY, data.Title, data.Description, common.PENDING, userID, pq.Array(convertIntoIntegerSlice(data.Approvers)))
	if err != nil{
		log.Println("Error in creating task : ", err)
		return 0, err
	}
	var lastInsertID uint64
	err = trm.db.QueryRow(RETURN_TASK_ID_QUERY, data.Title, data.Description, common.PENDING, userID, pq.Array(convertIntoIntegerSlice(data.Approvers))).Scan(&lastInsertID)
	if err != nil {
		log.Println("Error retrieving last insert ID: ", err)
		return 0, err
	}
	return lastInsertID, nil
}

func (trm *TasksRepositoryManager) GetTaskByTaskID(ctx context.Context, taskID uint64) (common.Task, error){
	var intApprovers []int64
	var task common.Task
	err := trm.db.QueryRow(GET_TASK_BY_TASKID_QUERY, taskID).Scan(&task.ID, &task.Title, &task.Description, &task.TaskStatus, &task.CreatorID, pq.Array(&intApprovers), &task.CreatedAt)
	if err != nil {
		log.Println("Error in fetching task : ", err)
		return common.Task{}, nil
	}
	task.Approvers = convertIntoUsignedSlice(intApprovers)
	return task, nil
}

func (trm *TasksRepositoryManager) GetTasksByUserID(ctx context.Context) ([]common.Task, error) {
	userID := ctx.Value(common.USER_ID_KEY)
	rows, err := trm.db.Query(GET_TASKS_BY_USERID_QUERY, userID)
	if err != nil {
		log.Println("Error in fetching tasks : ", err)
		return nil, err
	}
	return scanTasksRecords(rows)
}

func (trm *TasksRepositoryManager) UpdateTaskByTaskID(ctx context.Context, data common.Task) (common.Task, error) {
	_, err := trm.db.Exec(UPDATE_TASK_BY_TASKID, data.Title, data.Description, data.TaskStatus, data.CreatorID, pq.Array(convertIntoIntegerSlice(data.Approvers)), data.ID)
	if err != nil {
		log.Println("Error in updating task : ", err)
		return common.Task{}, err
	}
	return data, nil
}

func scanTasksRecords(rows *sql.Rows) ([]common.Task, error) {
	var tasks []common.Task
	for rows.Next() {
		var task common.Task
		var intApprovers []int64
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.TaskStatus, &task.CreatorID, pq.Array(&intApprovers), &task.CreatedAt)
		if err != nil {
			log.Println("Error in scanning tasks : ", err)
			return nil, err
		}
		task.Approvers = convertIntoUsignedSlice(intApprovers)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func convertIntoIntegerSlice(values []uint64) []int64 {
	var intValues []int64
	for _, val := range values {
		intValues = append(intValues, int64(val))
	}
	return intValues
}

func convertIntoUsignedSlice(values []int64) []uint64 {
	var uValues []uint64
	for _, val := range values {
		uValues = append(uValues, uint64(val))
	}
	return uValues
}