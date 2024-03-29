package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi"

	"restapi/internal/domain/models"
	"restapi/internal/domain/services"
)

type UserHandler struct {
	userService *services.UserService
}

type User struct {
    ID    int
    Name  string
    Email string
}

func NewUserHandler(dbConn *sql.DB) *UserHandler {
	userRepo := models.NewUserRepository(dbConn)
	userService := services.NewUserService(userRepo)

	return &UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", uh.GetAllUsers)
	r.Post("/", uh.CreateUser)

	return r
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.userService.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("D:/code/go-chi-restapi-jwtauth/template/users.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, users)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting all users: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	err := uh.userService.Create(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error post user: %v", err), http.StatusInternalServerError)
		sentry.CaptureException(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
