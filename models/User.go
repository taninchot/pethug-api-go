package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID      `db:"id" binding:"required"`
	UserName  string         `db:"user_name" binding:"required"`
	MobileNo  string         `db:"mobile_no" binding:"required"`
	UserImage sql.NullString `db:"user_image"`
	CreatedAt time.Time      `db:"created_at" binding:"required"`
	UpdatedAt time.Time      `db:"updated_at" binding:"required"`
}
