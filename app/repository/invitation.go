package repository

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/utils"

	"github.com/google/uuid"
)

type InvitationRepository struct {
}

func NewInvitationRepository() InvitationRepository {
	return InvitationRepository{}
}

// GetSessionIdByVoterId retrieves the session ID by voter ID
func (p InvitationRepository) GetSessionIdByVoterId(sessionrId uint64) (uuid.UUID, error) {
	var sessionId uuid.UUID
	err := database.SqlSession.Model(&model.Invitation{}).
		Select("session_id").
		Where("id = ?", sessionrId).
		Scan(&sessionId).Error

	if err != nil {
		return uuid.Nil, err
	}

	return sessionId, nil
}

// GetInvitations retrieves all invitations based on the provided conditions.
func (p InvitationRepository) GetInvitations(invitationQuery *model.InvitationQuery) ([]model.Invitation, int64, error) {
	var invitations []model.Invitation
	var total int64

	query := database.SqlSession.Model(&model.Invitation{})

	// SessionID Like filter
	if invitationQuery.SessionID != uuid.Nil {
		query = query.Where("session_id = ?", invitationQuery.SessionID)
	}

	// Status Exact match
	if invitationQuery.Status {
		query = query.Where("status = ?", invitationQuery.Status)
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Retrieve data
	// Use pagination service to handle pagination
	paginationRepository := NewPaginationRepository[*model.InvitationQuery, model.Invitation]()
	query, err = paginationRepository.Handler(query, invitationQuery)
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&invitations).Error; err != nil {
		return nil, 0, err
	}

	// Query data
	err = query.Find(&invitations).Error

	return invitations, total, err
}

// CreateInvitations creates invitations based on the provided session ID and password list.
func (p InvitationRepository) CreateInvitations(sessionId uuid.UUID, number uint, length uint, format string) ([]*model.Invitation, error) {
	invitationUtil := &utils.Invitation{}
	// Generate Invitations
	invitations, err := invitationUtil.GenerateInvitation(number, length, format)
	if err != nil {
		return nil, err
	}

	// Encrypt Invitations
	invitationModels := make([]*model.Invitation, len(invitations))
	for i, invitation := range invitations {
		invitationEncrypt, err := invitationUtil.Encrypt(invitation)
		if err != nil {
			return nil, err
		}
		invitationModels[i] = &model.Invitation{
			SessionID:  sessionId,
			CodeHash: invitationEncrypt,
		}
	}

	// Use transaction to ensure all invitations are created successfully
	transaction := database.SqlSession.Begin()
	err = transaction.CreateInBatches(&invitationModels, 100).Error

	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	return invitationModels, transaction.Commit().Error
}

// UpdateInvitationStatus updates the status of invitations based on the provided session ID and invitation ID list.
func (p InvitationRepository) UpdateInvitationStatus(sessionId uuid.UUID, invitationIDs []uint64, status bool) ([]*model.Invitation, error) {
	var invitationModels []*model.Invitation
	err := database.SqlSession.
		Model(&model.Invitation{}).
		Where("session_id = ? AND id IN ?", sessionId, invitationIDs).
		Update("status", status).
		Scan(&invitationModels).Error

	if err != nil {
		return nil, err
	}

	return invitationModels, nil
}

// DeleteInvitation deletes invitations based on the provided invitation ID list.
func (p InvitationRepository) DeleteInvitation(ids []uint64, userInfo model.UserInfo) ([]*model.Invitation, error) {
	var invitations []*model.Invitation

	query := database.SqlSession.
		Model(&model.Invitation{}).
		Where("invitations.id IN ?", ids)

	// Non-admin users need to check the ownership
	if !userInfo.IsAdmin {
		query = query.
			Joins("JOIN sessions ON invitations.session_id = sessions.uuid").
			Where("sessions.user_id = ?", userInfo.UserID)
	}

	// Retrieve invitations and delete them directly
	if err := query.Find(&invitations).Error; err != nil {
		return nil, err
	}

	if len(invitations) == 0 {
		return invitations, nil
	}

	// Extract IDs and delete
	invitationIds := make([]uint64, len(invitations))
	for i, invitation := range invitations {
		invitationIds[i] = invitation.ID
	}

	deleteError := database.SqlSession.Delete(&model.Invitation{}, invitationIds).Error
	return invitations, deleteError
}
