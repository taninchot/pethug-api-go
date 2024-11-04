package dtos

import (
	"github.com/google/uuid"
)

type UserCreateRes struct {
	Id          uuid.UUID `json:"id"`
	UserName    string    `json:"userName"`
	UserImage   *string   `json:"userImage"`
	MobileNo    string    `json:"mobileNo"`
	AccessToken string    `json:"accessToken"`
}
