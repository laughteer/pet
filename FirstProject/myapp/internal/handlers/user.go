package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"myapp/internal/models"

	"github.com/gorilla/mux"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		http.Error(w, "failed to insert user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, "failed to scan user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var user models.User
	err = h.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	res, err := h.db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, id)
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil || count == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	res, err := h.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil || count == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
