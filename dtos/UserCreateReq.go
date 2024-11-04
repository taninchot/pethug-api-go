package dtos

type UserCreateReq struct {
	UserName  string  `json:"userName" binding:"required"`
	MobileNo  string  `json:"mobileNo" binding:"required"`
	UserImage *string `json:"userImage"`
	Password  string  `json:"password" binding:"required"`
}
