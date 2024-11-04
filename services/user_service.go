package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"pethug-api-go/dtos"
	"pethug-api-go/models"
	"pethug-api-go/repositories"
	"pethug-api-go/utils"
	"time"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]dtos.UserGetListRes, error) {
	users, err := s.userRepository.GetAllUsers(ctx)
	return users, err
}

func (s *UserService) CreateUser(ctx context.Context, userCreateReq dtos.UserCreateReq) (dtos.UserCreateRes, error) {
	isUserExists, err := s.userRepository.CheckIsExistsByUserNameOrMobileNo(ctx, userCreateReq.MobileNo, userCreateReq.UserName)
	if err != nil {
		return dtos.UserCreateRes{}, err
	}
	if isUserExists {
		return dtos.UserCreateRes{}, fmt.Errorf("user with mobile number %s or username %s already exists", userCreateReq.MobileNo, userCreateReq.UserName)
	}

	// Prepare Request Nullable Fields
	var userImageReq sql.NullString
	if userCreateReq.UserImage != nil && *userCreateReq.UserImage != "" {
		userImageReq = sql.NullString{String: *userCreateReq.UserImage, Valid: true}
	} else {
		userImageReq = sql.NullString{Valid: false}
	}

	// Begin Transaction
	tx, err := s.userRepository.BeginTx(ctx)
	if err != nil {
		return dtos.UserCreateRes{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				log.Fatalf("Rolling Back %v\n", err)
				return
			}
			panic(p)
		} else if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				log.Fatalf("Rolling Back %v\n", err)
				return
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()

	createdUser, err := s.userRepository.CreateUserTx(ctx, tx, models.User{
		Id:        uuid.New(),
		UserName:  userCreateReq.UserName,
		MobileNo:  userCreateReq.MobileNo,
		UserImage: userImageReq,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return dtos.UserCreateRes{}, err
	}

	hashPassword, err := hashPassword(userCreateReq.Password)
	err = s.userRepository.CreatePasswordUserTx(ctx, tx, createdUser.Id, hashPassword)

	if err != nil {
		log.Fatalln("Cannot Create PasswordUser.")
	}

	// Prepared Response Nullable Fields
	var userImageResponse *string
	if createdUser.UserImage.Valid {
		userImageResponse = &createdUser.UserImage.String
	} else {
		userImageResponse = nil
	}

	accessToken, err := utils.GenerateJWT(createdUser.Id)
	if err != nil {
		log.Fatalln("Cannot Create AccessToken.")
	}

	return dtos.UserCreateRes{
		Id:          createdUser.Id,
		UserName:    createdUser.UserName,
		MobileNo:    createdUser.MobileNo,
		UserImage:   userImageResponse,
		AccessToken: accessToken,
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
