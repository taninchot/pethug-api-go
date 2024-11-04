package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"pethug-api-go/dtos"
	"pethug-api-go/services"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (u *UserController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/users/get/list", u.GetUsers)
	router.POST("/users/create", u.CreateUser)
}

func (u *UserController) GetUsers(ctx *gin.Context) {
	users, err := u.service.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	var userRequest dtos.UserCreateReq
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		var errorMessages []string

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldErr := range validationErrors {
				errorMessages = append(errorMessages, fieldErr.Error())
			}
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"reason": errorMessages})
		return
	}

	userCreated, err := u.service.CreateUser(ctx, userRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userCreated)
}

func (u UserController) LoginUser() {

}
