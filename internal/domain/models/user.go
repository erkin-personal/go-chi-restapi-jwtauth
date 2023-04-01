package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetAll() ([]*User, error) {
	rows, err := ur.db.Query("SELECT id, username, email, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) Create(user *User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := ur.db.QueryRow("INSERT INTO users(username, email, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING id",
		user.Username, user.Email, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)

	return err
}
