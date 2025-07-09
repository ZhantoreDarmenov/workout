package services

import (
	"context"

	"workout/internal/models"
	"workout/internal/repositories"
)

// InviteService handles business logic for program invitations.
type InviteService struct {
	Repo     *repositories.InviteRepository
	UserRepo *repositories.UserRepository
}

func (s *InviteService) InviteClient(ctx context.Context, programID int, email, message string, days int) (models.ProgramInvite, error) {
	invite := models.ProgramInvite{
		ProgramID:  programID,
		Email:      email,
		Message:    message,
		AccessDays: days,
	}
	return s.Repo.CreateInvite(ctx, invite)
}

func (s *InviteService) AcceptInvite(ctx context.Context, token string, clientID int) (models.ProgramInvite, error) {
	return s.Repo.AcceptInvite(ctx, token, clientID)
}

func (s *InviteService) UpdateAccess(ctx context.Context, programID, clientID, days int) (models.ProgramInvite, error) {
	return s.Repo.UpdateAccessDuration(ctx, programID, clientID, days)
}
