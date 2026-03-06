package service

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/repository"
	"vote/app/utils"

	"github.com/google/uuid"
)

type InvitationService struct {
}

func NewInvitationService() InvitationService {
	return InvitationService{}
}

// GetInvitation retrieves an invitation based on the provided vote ID and invitation string.
func (p InvitationService) GetInvitation(voteId uuid.UUID, invitation string) (*model.Invitation, error) {
	invitationModel := model.Invitation{}
	err := database.SqlSession.
		Where("vote_id = ? AND invitation = ? AND status = true", voteId, invitation).
		First(&invitationModel).
		Error

	if err != nil {
		return nil, err
	}

	return &invitationModel, nil
}

// GetInvitations retrieves all invitations based on the provided vote ID and query parameters.
func (p InvitationService) GetInvitations(voteId uuid.UUID, invitationQuery *model.InvitationQuery) ([]*model.InvitationConnection, error) {
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
func (p InvitationService) CreateInvitation(voteId uuid.UUID, number uint, length uint, format string) ([]*model.Invitation, error) {
	invitations, err := repository.NewInvitationRepository().CreateInvitations(voteId, number, length, format)
	if err != nil {
		return nil, err
	}

	return invitations, nil
}

// UpdateInvitationStatus updates the status of invitations based on the provided vote ID and invitation IDs.
func (p InvitationService) UpdateInvitationStatus(voteId uuid.UUID, invitationIDs []uint64, status bool) ([]*model.Invitation, error) {
	invitations, err := repository.NewInvitationRepository().UpdateInvitationStatus(voteId, invitationIDs, status)
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
