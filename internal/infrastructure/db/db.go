package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	_ "github.com/lib/pq"

)

type User struct {
    ID    int
    Name  string
    Email string
}

func NewConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err != nil {
		sentry.CaptureException(err)
		// handle the error
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getAllUsers(db *sql.DB) ([]User, error) {
    rows, err := db.Query("SELECT id, name, email FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}

