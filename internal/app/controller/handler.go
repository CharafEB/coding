package controller

import (
	"database/sql"
	"encoding/json"
	"github/codingMaster/internal/app/model"
	"log"
	"net/http"
	"strconv"
)

func (t *Application) tst(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

type articlesC struct {
	DB *sql.DB
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	var login model.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if _, err := app.Store.Admin.Login(r.Context(), &login); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func (app *Application) Singup(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := app.Store.Admin.Signup(r.Context(), &user); err != nil {
		log.Print(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// Add Project Handler
func (app *Application) AddProject(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.AddProject(r.Context(), &project); err != nil {
		log.Printf("Error adding project: %v", err)
		http.Error(w, "Failed to add project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Project added successfully"))
}

// Delete Project Handler
func (app *Application) DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.DeleteProject(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project deleted successfully"))
}

// Update Project Handler
func (app *Application) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.UpdateProject(r.Context(), &project); err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project updated successfully"))
}

// Upprove Handler
func (app *Application) Upprove(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.Upprove(r.Context(), id, payload.Status); err != nil {
		http.Error(w, "Failed to update project status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project status updated successfully"))
}

// GetProject Handler
func (app *Application) GetProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := app.Store.Role.GetProject(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (app *Application) Unpprove(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.Unpprove(r.Context(), id); err != nil {
		http.Error(w, "Failed to reject project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project rejected successfully"))
}
