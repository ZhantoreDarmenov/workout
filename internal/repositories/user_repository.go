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

func (r *UserRepository) GetUserByPhone(ctx context.Context, phone string) (models.User, error) {
	var user models.User
	query := `
        SELECT id, name, phone, email, password, role, created_at, updated_at
        FROM users
        WHERE phone = ?
    `
	err := r.DB.QueryRowContext(ctx, query, phone).Scan(
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

func (r *UserRepository) GetVerificationCodeByPhone(ctx context.Context, phone string) (string, error) {
	var code string
	err := r.DB.QueryRowContext(ctx, `SELECT code FROM verification_codes WHERE phone = ?`, phone).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", models.ErrInvalidVerificationCode
		}
		return "", err
	}
	return code, nil
}

func (r *UserRepository) ClearVerificationCode(ctx context.Context, phone string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM verification_codes WHERE phone = ?`, phone)
	return err
}
