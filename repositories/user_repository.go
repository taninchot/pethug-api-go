package repositories

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pethug-api-go/db"
	"pethug-api-go/dtos"
	"pethug-api-go/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	// Acquire a connection from the pool
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}

	// Begin a transaction on the acquired connection
	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release() // Ensure to release the connection back to the pool if Begin fails
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Return the transaction object; remember to release connection when transaction ends
	return tx, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]dtos.UserGetListRes, error) {
	rows, err := db.DB.Query(ctx, "SELECT id, user_name, mobile_no FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]dtos.UserGetListRes, 0)
	for rows.Next() {
		var user dtos.UserGetListRes
		if err := rows.Scan(&user.Id, &user.UserName, &user.MobileNo); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CheckIsExistsByUserNameOrMobileNo(ctx context.Context, mobileNo string, userName string) (bool, error) {
	rows, err := db.DB.Query(ctx, "SELECT 1 FROM users WHERE mobile_no = $1 OR LOWER(user_name) = LOWER($2)", mobileNo, userName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (r *UserRepository) CreateUserTx(ctx context.Context, tx pgx.Tx, user models.User) (models.User, error) {
	_, err := tx.Exec(ctx, "INSERT INTO users (id, user_name, mobile_no, user_image, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Id, user.UserName, user.MobileNo, user.UserImage, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) CreatePasswordUserTx(ctx context.Context, tx pgx.Tx, userId uuid.UUID, hashedPassword string) error {
	_, err := tx.Exec(ctx, "INSERT INTO password_users(id, user_id, hash_password) VALUES ($1, $2, $3)",
		uuid.New(), userId, hashedPassword)
	return err
}
