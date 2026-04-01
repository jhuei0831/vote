package service

import (
	"context"
	"fmt"
	"vote/app/database"
	"vote/app/middleware"
	"vote/app/model"
	"vote/app/repository"
	"vote/app/utils"
	pb "vote/proto/voter"

	"github.com/google/uuid"
)

type InvitationService struct {
}

func NewInvitationService() InvitationService {
	return InvitationService{}
}

// GetInvitation retrieves an invitation based on the provided vote ID and invitation string.
func (p InvitationService) GetInvitation(sessionId uuid.UUID, invitation string) (*model.Invitation, error) {
	invitationModel := model.Invitation{}
	err := database.SqlSession.
		Where("session_id = ? AND code_hash = ? AND status = true", sessionId, invitation).
		Preload("Ballots").
		First(&invitationModel).
		Error

	if err != nil {
		return nil, err
	}

	return &invitationModel, nil
}

// GetInvitations retrieves all invitations based on the provided vote ID and query parameters.
func (p InvitationService) GetInvitations(sessionId uuid.UUID, invitationQuery *model.InvitationQuery) ([]*model.InvitationConnection, error) {
	invitations, total, err := repository.NewInvitationRepository().GetInvitations(invitationQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.InvitationQuery, model.Invitation]()
	invitations, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(invitations, invitationQuery)

	paginationService := NewPaginationService[model.Invitation, model.InvitationEdge, *model.InvitationConnection]()
	connection := paginationService.BuildConnection(invitations, total, hasPreviousPage, hasNextPage,
		func(invitation model.Invitation) uint64 {
			return invitation.ID
		},
	)

	return []*model.InvitationConnection{connection}, nil
}

// DecryptInvitation decrypts the provided invitation string and returns the original invitation.
func (p InvitationService) DecryptInvitation(invitation string) (string, error) {
	decryptedInvitation, err := (&utils.Invitation{}).Decrypt(invitation)
	if err != nil {
		return "", err
	}

	return decryptedInvitation, nil
}

// CreateInvitation Create invitations can encrypt and decrypt
func (p InvitationService) CreateInvitation(sessionId uuid.UUID, number uint, length uint, format string) ([]*model.Invitation, error) {
	invitations, err := repository.NewInvitationRepository().CreateInvitations(sessionId, number, length, format)
	if err != nil {
		return nil, err
	}

	return invitations, nil
}

// UpdateInvitationStatus updates the status of invitations based on the provided vote ID and invitation IDs.
func (p InvitationService) UpdateInvitationStatus(sessionId uuid.UUID, invitationIDs []uint64, status bool) ([]*model.Invitation, error) {
	invitations, err := repository.NewInvitationRepository().UpdateInvitationStatus(sessionId, invitationIDs, status)
	if err != nil {
		return nil, err
	}

	return invitations, nil
}

// DeleteInvitation deletes invitations based on the provided invitation IDs and user information.
func (p InvitationService) DeleteInvitation(ids []uint64, userInfo model.UserInfo) ([]*model.Invitation, error) {
	invitations, err := repository.NewInvitationRepository().DeleteInvitation(ids, userInfo)
	if err != nil {
		return nil, err
	}

	return invitations, nil
}

// VerifyInviteCode verifies an invite code using the gRPC VoterService.
func (p InvitationService) VerifyInviteCode(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	code := req.GetCode()
	sessionIDStr := req.GetSessionId()

	if code == "" {
		return &pb.ValidateResponse{
			Success: false,
			Message: "Invite code cannot be empty",
		}, nil
	}

	hashCode, err := (&utils.Invitation{}).Encrypt(code)
	if err != nil {
		return &pb.ValidateResponse{
			Success: false,
			Message: "Failed to encrypt invite code",
		}, nil
	}

	// Check invitation exist
	sessionId, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return &pb.ValidateResponse{
			Success: false,
			Message: "Invalid session ID format",
		}, nil
	}

	invitation, err := p.GetInvitation(sessionId, hashCode)
	if err != nil {
		return &pb.ValidateResponse{
			Success: false,
			Message: fmt.Sprintf("failed to validate invitation: %v", err),
		}, nil
	}

	// Record invitation usage
	invitationUsage := model.InvitationUsage{
		InvitationID: invitation.ID,
		VoterTempID:  uuid.New(),
	}

	insertErr := database.SqlSession.Model(&model.InvitationUsage{}).Create(&invitationUsage).Error
	if insertErr != nil {
		return &pb.ValidateResponse{
			Success: false,
			Message: fmt.Sprintf("failed to record invitation usage: %v", insertErr),
		}, nil
	}

	isVoted := len(invitation.Ballots) > 0
	tokenString, _, err := middleware.GenVoterToken(invitationUsage, sessionId, isVoted)
	if err != nil {
		return &pb.ValidateResponse{
			Success: false,
			Message: fmt.Sprintf("failed to generate JWT token: %v", err),
		}, nil
	}

	return &pb.ValidateResponse{
		Success: true,
		Jwt:     tokenString,
		Message: "Invite code verified successfully",
	}, nil
}
