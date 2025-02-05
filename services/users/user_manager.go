package users

import (
	"context"
	"log"
	"net/mail"
	"task-manager-backend/data/common"
	"task-manager-backend/data/contracts"
	"task-manager-backend/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserServiceManager struct{
	userRepo 		types.UsersRepository
	secret 			string
}

func NewUserServiceManager(userDb types.UsersRepository, secret string) *UserServiceManager {
	return &UserServiceManager{
		userRepo: userDb,
		secret: secret,
	}
}

func (usm *UserServiceManager) SignUp(ctx context.Context, r contracts.SignUpRequest) contracts.SignUpResponse {
	err := validateUserData(r.UserData)
	if err != nil {
		return contracts.SignUpResponse{
			Status: "FAILED : Invalid email address [" + err.Error() + "]",
		}
	} 
	err = usm.userRepo.CreateUser(ctx, r.UserData)
	if err != nil {
		return contracts.SignUpResponse{
			Status: "FAILED : " + err.Error(),
		}
	}
	return contracts.SignUpResponse{
		Status: "SUCCESS",
	}
}

func (usm *UserServiceManager) Login(ctx context.Context, r contracts.LoginRequest) contracts.LoginResponse {
	user, err := usm.userRepo.GetUserByName(ctx, r.Name)
	if err != nil {
		return contracts.LoginResponse{
			Status: "FAILED : " + err.Error(),
		}
	}
	
	// Verify typed password
	if r.Password != user.Password {
		return contracts.LoginResponse{
			Status: "FAILED : Password donot match !",
		}
	}
	
	// Generate a new JWT token for the user valid for 5 minutes
	expirationTime := time.Now().Add(common.TOKEN_EXPIRY_MINUTES * time.Minute)
	claims := &common.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Use the signing method and the secret to sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(usm.secret))
	if err != nil {
		return contracts.LoginResponse{
			Status: "FAILED : " + err.Error(),
		}
	}
	return contracts.LoginResponse{
		Status: "SUCCESS",
		Token: tokenString,
	}
}

func (usm *UserServiceManager) GetAllUsers(ctx context.Context) contracts.GetAllUsersResponse {
	users, err := usm.userRepo.GetAllUsers(ctx)
	if err != nil {
		return contracts.GetAllUsersResponse{
			Status: "FAILED : " + err.Error(),
		}
	}
	return contracts.GetAllUsersResponse{
		Status: "SUCCESS",
		Users: users,
	}
}

func validateUserData(r common.UserData) error {
	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		log.Println("Error in parsing email : ", err)
		return err
	}
	return nil
}