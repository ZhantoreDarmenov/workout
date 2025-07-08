package repositories

import (
	"context"
	"database/sql"
	"errors"
	_ "fmt"
	_ "strings"
	"time"
	_ "time"
	"workout/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	DB *sql.DB
}

type Session struct {
	ID     string `json:"id"`
	Expiry string `json:"expiry"`
}

func (r *UserRepository) SetSession(ctx context.Context, id string, session models.Session) error {

	query := `
		UPDATE users 
		SET refresh_token = ?, expires_at = ? 
		WHERE id = ?
	`

	result, err := r.DB.ExecContext(ctx, query, session.RefreshToken, session.ExpiresAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (r *UserRepository) GetSession(ctx context.Context, id string) (models.Session, error) {
	query := `
		SELECT refresh_token, expires_at
		FROM users
		WHERE id = ?
	`

	var session models.Session
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return session, errors.New("no session found for the user")
		}
		return session, err
	}

	return session, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := `
        SELECT id, name, phone, email, password, role, created_at, updated_at
        FROM users
        WHERE email = ?
    `
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Phone, &user.Email, &user.Password,
		&user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `
        INSERT INTO users (name, phone, email, password, role, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	user.CreatedAt = time.Now()
	user.UpdatedAt = &user.CreatedAt
	result, err := r.DB.ExecContext(ctx, query,
		user.Name, user.Phone, user.Email, user.Password, user.Role,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}
	user.ID = int(id)
	return user, nil
}

func (r *UserRepository) GetVerificationCodeByEmail(ctx context.Context, email string) (string, error) {
	var code string
	err := r.DB.QueryRowContext(ctx, `SELECT code FROM verification_codes WHERE email = ?`, email).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", models.ErrInvalidVerificationCode
		}
		return "", err
	}
	return code, nil
}

// ClearVerificationCode removes a verification code record for an email.
func (r *UserRepository) ClearVerificationCode(ctx context.Context, email string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM verification_codes WHERE email = ?`, email)
	return err
}

// UpdateUserRole updates the role of a user.
func (r *UserRepository) UpdateUserRole(ctx context.Context, userID int, role string) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE users SET role = ?, updated_at = ? WHERE id = ?`, role, time.Now(), userID)
	return err
}

func (r *UserRepository) GetAllClients(ctx context.Context) ([]models.User, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, phone, email, password, role, created_at, updated_at FROM users WHERE role = 'client'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

func (r *UserRepository) GetClientsByProgramID(ctx context.Context, programID int) ([]models.User, error) {
	query := `SELECT DISTINCT u.id, u.name, u.phone, u.email, u.password, u.role, u.created_at, u.updated_at
              FROM users u
              JOIN progress p ON u.id = p.client_id
              JOIN days d ON p.day_id = d.id
              WHERE d.work_out_program_id = ?`
	rows, err := r.DB.QueryContext(ctx, query, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

func (r *UserRepository) DeleteClientFromProgram(ctx context.Context, programID, clientID int) error {
	query := `DELETE p FROM progress p JOIN days d ON p.day_id = d.id WHERE d.work_out_program_id = ? AND p.client_id = ?`
	_, err := r.DB.ExecContext(ctx, query, programID, clientID)
	return err
}
func (r *UserRepository) GetProgramsByClientID(ctx context.Context, clientID int) ([]models.WorkOutProgram, error) {
	query := `SELECT DISTINCT wp.id, wp.trainer_id, wp.name, wp.days, wp.description, wp.created_at, wp.updated_at
                 FROM workout_programs wp
                 JOIN days d ON wp.id = d.work_out_program_id
                 JOIN progress p ON d.id = p.day_id
                 WHERE p.client_id = ?`
	rows, err := r.DB.QueryContext(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.WorkOutProgram
	for rows.Next() {
		var p models.WorkOutProgram
		if err := rows.Scan(&p.ID, &p.TrainerID, &p.Name, &p.Days, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}
