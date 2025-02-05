package psql

import (
	"context"
	"database/sql"
	"log"
	"task-manager-backend/data/common"
)

const (
	CREAT_USER_QUERY = "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	GET_USER_BY_USERID_QUERY = "SELECT * FROM users WHERE id = $1"
	GET_USER_BY_NAME_QUERY = "SELECT * FROM users WHERE name = $1"
	GET_ALL_USERS = "SELECT * FROM users"
)

type UsersRepositoryManager struct{
	db 		*sql.DB
}

func NewUsersRepositoryManager(client *sql.DB) *UsersRepositoryManager{
	return &UsersRepositoryManager{
		db: client,
	}
}

func (urm *UsersRepositoryManager) CreateUser(ctx context.Context, data common.UserData) error {
	_, err := urm.db.Exec(CREAT_USER_QUERY, data.Name, data.Email, data.Password)
	if err != nil {
		log.Println("Error while creating user : ", err)
		return err
	}
	return nil
}

func (urm *UsersRepositoryManager) GetUserByUserID(ctx context.Context, userID uint64) (common.User, error) {
	var user common.User
	err := urm.db.QueryRow(GET_USER_BY_USERID_QUERY, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		log.Println("Error while fetching user : ", err)
		return common.User{}, err
	}
	return user, nil
}

func (urm *UsersRepositoryManager) GetUserByName(ctx context.Context, name string) (common.User, error) {
	var user common.User
	err := urm.db.QueryRow(GET_USER_BY_NAME_QUERY, name).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		log.Println("Error while fetching user : ", err)
		return common.User{}, err
	}
	return user, nil
}

func (urm *UsersRepositoryManager) GetAllUsers(ctx context.Context) ([]common.User, error) {
	rows, err := urm.db.Query(GET_ALL_USERS)
	if err != nil {
		log.Println("Error while fecthing users : ", err)
		return nil, err
	}

	var users []common.User
	for rows.Next() {
		var user common.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			log.Println("Error in scanning users : ", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
