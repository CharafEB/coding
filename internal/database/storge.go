package database

import (
	"context"
	"database/sql"
	"fmt"
	"github/codingMaster/internal/app/model"
)

type articlesC struct {
	DB *sql.DB
}

func (a *articlesC) Login(ctx context.Context, login *model.Login) (*model.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	var user model.User
	err := a.DB.QueryRowContext(ctx, query, login.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", login.Email)
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	if user.Password != login.Password {
		return nil, fmt.Errorf("invalid password for user with email %s", login.Email)
	}

	return &user, nil
}

func (a *articlesC) Signup(ctx context.Context, val *model.User) error {
	// Check if the user already exists
	var exists bool
	err := a.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", val.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	if exists {
		return fmt.Errorf("user with email %s already exists", val.Email)
	}

	// Insert the new user
	query := `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`
	_, err = a.DB.ExecContext(ctx, query, val.ID, val.Name, val.Email, val.Password)
	if err != nil {
		return fmt.Errorf("error inserting new user: %w", err)
	}

	return nil // Signup successful
}

func (a *articlesC) AddProject(ctx context.Context, project *model.Project) error {
	_, err := a.DB.ExecContext(ctx, `
		INSERT INTO projects (title, student_id, project_partner_id, status, date, date_edit, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		project.Title, project.StidentIDB, project.ProjectPartnerIDs, project.Status, project.Date, project.DateEdit, project.SupportStucter)
	if err != nil {
		return fmt.Errorf("error adding project: %w", err)
	}
	return nil
}

func (a *articlesC) DeleteProject(ctx context.Context, id int) error {
	_, err := a.DB.ExecContext(ctx, "DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting project: %w", err)
	}
	return nil
}

// Update Project
func (a *articlesC) UpdateProject(ctx context.Context, project *model.Project) error {
	_, err := a.DB.ExecContext(ctx, `
		UPDATE projects
		SET title = $1, student_id = $2, project_partner_id = $3, status = $4, date = $5, date_edit = $6, role = $7
		WHERE id = $8`,
		project.Title, project.StidentIDB, project.ProjectPartnerIDs, project.Status, project.Date, project.DateEdit, project.SupportStucter, project.ID)
	if err != nil {
		return fmt.Errorf("error updating project: %w", err)
	}
	return nil
}

func (a *articlesC) Upprove(ctx context.Context, id int, status string) error {
	_, err := a.DB.ExecContext(ctx, "UPDATE projects SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return fmt.Errorf("error updating project status: %w", err)
	}
	return nil
}

// GetProject function to retrieve a project by ID
func (a *articlesC) GetProject(ctx context.Context, id int) (*model.Project, error) {
	project := &model.Project{}
	err := a.DB.QueryRowContext(ctx, `
		SELECT id, title, student_id, project_partner_id, status, date, date_edit, role
		FROM projects WHERE id = $1`, id).Scan(
		&project.ID,
		&project.Title,
		&project.StidentIDB,
		&project.ProjectPartnerIDs,
		&project.Status,
		&project.Date,
		&project.DateEdit,
		&project.SupportStucter,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("error retrieving project: %w", err)
	}
	return project, nil
}

func (a *articlesC) Unpprove(ctx context.Context, id int) error {
	_, err := a.DB.ExecContext(ctx, "UPDATE projects SET status = $1 WHERE id = $2", "rejected", id)
	if err != nil {
		return fmt.Errorf("error rejecting project: %w", err)
	}
	return nil
}
