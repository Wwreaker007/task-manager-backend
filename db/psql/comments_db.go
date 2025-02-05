package psql

import (
	"context"
	"database/sql"
	"log"
	"task-manager-backend/data/common"
)

const (
	CREATE_COMMENT_QUERY = "INSERT INTO task_comments (task_id, commentor, comment) VALUES ($1, $2, $3)"
	GET_COMMENTS_BY_TASKID_QUERY = "SELECT * FROM task_comments WHERE task_id = $1"
	GET_COMMENTS_BY_USERID_QUERY = "SELECT * FROM task_comments WHERE commentor = $1"
)

type CommentsRepositoryManager struct{
	db		*sql.DB
}

func NewCommentsRepositoryManager(client *sql.DB) *CommentsRepositoryManager{
	return &CommentsRepositoryManager{
		db: client,
	}
}

func (crm *CommentsRepositoryManager) CreateComment(ctx context.Context, data common.TaskCommentData) error {
	userID := ctx.Value(common.USER_ID_KEY)
	_, err := crm.db.Exec(CREATE_COMMENT_QUERY, data.TaskID, userID, data.Comment)
	if err != nil {
		log.Println("Error in creating comment : ", err)
		return err
	}
	return nil
}

func (crm *CommentsRepositoryManager) GetCommentsByTaskID(ctx context.Context, taskID uint64) ([]common.TaskComment, error) {
	rows, err := crm.db.Query(GET_COMMENTS_BY_TASKID_QUERY, taskID)
	if err != nil {
		log.Println("Error in fetching comments : ", err)
		return nil, err
	}
	return scanCommentRecords(rows)
}

func (crm *CommentsRepositoryManager) GetCommentsByUserID(ctx context.Context) ([]common.TaskComment, error) {
	userID := ctx.Value(common.USER_ID_KEY)
	rows, err := crm.db.Query(GET_COMMENTS_BY_TASKID_QUERY, userID)
	if err != nil {
		log.Println("Error in fetching comments : ", err)
		return nil, err
	}
	return scanCommentRecords(rows)
}

func scanCommentRecords(rows *sql.Rows) ([]common.TaskComment, error) {
	var comments []common.TaskComment
	for rows.Next() {
		var comment common.TaskComment
		err := rows.Scan(&comment.ID, &comment.TaskID, &comment.Commentor, &comment.Comment, &comment.CommentedAt)
		if err != nil {
			log.Println("Error in scanning comments : ", err)
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}