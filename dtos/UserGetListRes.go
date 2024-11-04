package dtos

import (
	"github.com/google/uuid"
)

type UserGetListRes struct {
	Id       uuid.UUID `json:"id"`
	UserName string    `json:"user_name"`
	MobileNo string    `json:"mobile_no"`
}
