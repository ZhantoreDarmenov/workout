package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"workout/internal/models"
)

// InviteRepository manages program invitation records.
type InviteRepository struct {
	DB *sql.DB
}

func (r *InviteRepository) CreateInvite(ctx context.Context, inv models.ProgramInvite) (models.ProgramInvite, error) {
	inv.CreatedAt = time.Now()
	token := uuid.New().String()
	inv.Token = token
	res, err := r.DB.ExecContext(ctx, `INSERT INTO program_invites (program_id, email, message, access_days, token, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		inv.ProgramID, inv.Email, inv.Message, inv.AccessDays, inv.Token, inv.CreatedAt)
	if err != nil {
		return models.ProgramInvite{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.ProgramInvite{}, err
	}
	inv.ID = int(id)
	return inv, nil
}

func (r *InviteRepository) getInviteByToken(ctx context.Context, token string) (models.ProgramInvite, error) {
	var inv models.ProgramInvite
	var clientID sql.NullInt64
	var acceptedAt, expires sql.NullTime
	query := `SELECT id, program_id, email, message, access_days, token, client_id, accepted_at, access_expires, created_at, updated_at FROM program_invites WHERE token = ?`
	err := r.DB.QueryRowContext(ctx, query, token).Scan(&inv.ID, &inv.ProgramID, &inv.Email, &inv.Message, &inv.AccessDays,
		&inv.Token, &clientID, &acceptedAt, &expires, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ProgramInvite{}, models.ErrInviteNotFound
		}
		return models.ProgramInvite{}, err
	}
	if clientID.Valid {
		cid := int(clientID.Int64)
		inv.ClientID = &cid
	}
	if acceptedAt.Valid {
		inv.AcceptedAt = &acceptedAt.Time
	}
	if expires.Valid {
		inv.AccessExpires = &expires.Time
	}
	return inv, nil
}

func (r *InviteRepository) AcceptInvite(ctx context.Context, token string, clientID int) (models.ProgramInvite, error) {
	inv, err := r.getInviteByToken(ctx, token)
	if err != nil {
		return models.ProgramInvite{}, err
	}
	now := time.Now()
	expires := now.Add(time.Duration(inv.AccessDays) * 24 * time.Hour)
	_, err = r.DB.ExecContext(ctx, `UPDATE program_invites SET client_id=?, accepted_at=?, access_expires=?, updated_at=? WHERE id=?`,
		clientID, now, expires, now, inv.ID)
	if err != nil {
		return models.ProgramInvite{}, err
	}
	inv.ClientID = &clientID
	inv.AcceptedAt = &now
	inv.AccessExpires = &expires
	inv.UpdatedAt = &now
	return inv, nil
}

func (r *InviteRepository) UpdateAccessDuration(ctx context.Context, programID, clientID, days int) (models.ProgramInvite, error) {
	var inv models.ProgramInvite
	var acceptedAt sql.NullTime
	query := `SELECT id, accepted_at FROM program_invites WHERE program_id = ? AND client_id = ?`
	err := r.DB.QueryRowContext(ctx, query, programID, clientID).Scan(&inv.ID, &acceptedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ProgramInvite{}, models.ErrInviteNotFound
		}
		return models.ProgramInvite{}, err
	}
	if acceptedAt.Valid {
		now := time.Now()
		expires := acceptedAt.Time.Add(time.Duration(days) * 24 * time.Hour)
		_, err = r.DB.ExecContext(ctx, `UPDATE program_invites SET access_days=?, access_expires=?, updated_at=? WHERE id=?`,
			days, expires, now, inv.ID)
		if err != nil {
			return models.ProgramInvite{}, err
		}
		inv.AccessDays = days
		inv.AccessExpires = &expires
		inv.AcceptedAt = &acceptedAt.Time
		inv.ProgramID = programID
		cid := clientID
		inv.ClientID = &cid
		inv.ID = inv.ID
		inv.UpdatedAt = &now
		return inv, nil
	}
	// if not accepted yet just update access_days
	_, err = r.DB.ExecContext(ctx, `UPDATE program_invites SET access_days=? WHERE id=?`, days, inv.ID)
	if err != nil {
		return models.ProgramInvite{}, err
	}
	inv.AccessDays = days
	inv.ProgramID = programID
	cid := clientID
	inv.ClientID = &cid
	return inv, nil
}


func (r *InviteRepository) GetProgramFromInvite(ctx context.Context, token string) (models.WorkOutProgram, error) {
	var p models.WorkOutProgram
	query := `SELECT wp.id, wp.trainer_id, wp.name, wp.days, wp.description, wp.created_at, wp.updated_at
              FROM workout_programs wp
              JOIN program_invites pi ON pi.program_id = wp.id
              WHERE pi.token = ?`
	err := r.DB.QueryRowContext(ctx, query, token).Scan(&p.ID, &p.TrainerID, &p.Name, &p.Days, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.WorkOutProgram{}, models.ErrInviteNotFound
		}
		return models.WorkOutProgram{}, err
	}
	return p, nil
}

