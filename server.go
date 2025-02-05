package main

import (
	"database/sql"
	"net/http"
	"task-manager-backend/db/psql"
	"task-manager-backend/handler"
	"task-manager-backend/middleware"
	"task-manager-backend/services/notification"
	"task-manager-backend/services/tasks"
	"task-manager-backend/services/users"

	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
)

const (
	// User handlers API routes
	SIGNUP_USER 	= "/api/user/signup"
	LOGIN_USER  	= "/api/user/login"
	GET_USERS		= "/api/user/get"

	// Task handlers API routes
	CREATE_TASK 	= "/api/task/create"
	APPROVE_TASK 	= "/api/task/approve"
	COMMENT_TASK 	= "/api/task/comment"
	GET_TASK 		= "/api/task/get/{id}"
)

type Server struct {
	path 			string
	port 			string
	secret			string
	dbClient		*sql.DB
	smtpDialer		*gomail.Dialer
}

func NewServer(path string, port string, secret string, dbClient *sql.DB, dialer *gomail.Dialer) *Server{
	return &Server{
		path:	path,
		port: 	port,
		dbClient: dbClient,
		secret: secret,
		smtpDialer: dialer,
	}
}

func (s *Server) ServerStartup() error{
	serverMux := mux.NewRouter()

	// Assign the dependencies here which will be used by the services
	userRepo := psql.NewUsersRepositoryManager(s.dbClient)
	taskRepo := psql.NewTasksRepositoryManager(s.dbClient)
	approvalRepo := psql.NewApprovalsRepositoryManager(s.dbClient)
	commentsRepo := psql.NewCommentsRepositoryManager(s.dbClient)

	// Services which will be used by the handlers
	emailNotification := notification.NewEmailNotificationServiceManager(s.smtpDialer)
	taskService := tasks.NewTaskServiceManager(taskRepo, commentsRepo, approvalRepo, userRepo, emailNotification)
	userService := users.NewUserServiceManager(userRepo, s.secret)

	// Assign the middlewares which are requiered for API protection
	jwtAuth := middleware.NewJWTAuthorizationMiddleware(s.secret)

	// Handle routes for Signup and Login
	serverMux.HandleFunc(SIGNUP_USER, handler.NewSignUpHandler(userService).Handle).Methods("POST")
	serverMux.HandleFunc(LOGIN_USER, handler.NewLoginHandler(userService).Handle).Methods("POST")

	// Handle authentication of the routes
	serverMux.HandleFunc(GET_USERS, jwtAuth.Authenticate(handler.NewGetUsersHandler(userService).Handle)).Methods("GET")
	serverMux.HandleFunc(GET_TASK, jwtAuth.Authenticate(handler.NewGetTaskHandler(taskService).Handle)).Methods("GET")
	serverMux.HandleFunc(CREATE_TASK, jwtAuth.Authenticate(handler.NewCreateTaskHandler(taskService).Handle)).Methods("POST")
	serverMux.HandleFunc(APPROVE_TASK, jwtAuth.Authenticate(handler.NewApproveTaskHandler(taskService).Handle)).Methods("POST")
	serverMux.HandleFunc(COMMENT_TASK, jwtAuth.Authenticate(handler.NewCommentTaskHandler(taskService).Handle)).Methods("POST")

	return http.ListenAndServe(s.path + s.port, serverMux)
}