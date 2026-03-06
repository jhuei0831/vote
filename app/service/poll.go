package service

import (
	"fmt"
	"vote/app/model"
	"vote/app/repository"

	"github.com/google/uuid"
)

type PollService struct {
}

func NewPollService() PollService {
	return PollService{}
}

// GetPoll retrieves a poll by the provided ID and checks if it exists.
func (q PollService) GetPoll(id uint64, isAdmin bool, userId uint64) (*model.Poll, error) {
	return repository.NewPollRepository().GetPoll(id, isAdmin, userId, false)
}

// GetPollList retrieves poll options related to the provided SessionID.
func (q PollService) GetPollList(sessionID uuid.UUID) (*model.PollList, error) {
	polls, err := repository.NewPollRepository().GetPollsBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	return &model.PollList{
		List: polls,
	}, nil
}

// SelectPollWithPollOptions retrieves a poll by the provided ID and checks if it exists, along with its poll options.
func (q PollService) SelectPollWithPollOptions(id uint64, isAdmin bool, userId uint64) (*model.Poll, error) {
	return repository.NewPollRepository().GetPoll(id, isAdmin, userId, true)
}

// GetPolls retrieves polls based on the provided query parameters.
func (q PollService) GetPolls(pollQuery *model.PollQuery) ([]*model.PollConnection, error) {
	polls, total, err := repository.NewPollRepository().GetPolls(pollQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.PollQuery, model.Poll]()
	polls, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(polls, pollQuery)

	paginationService := NewPaginationService[model.Poll, model.PollEdge, *model.PollConnection]()
	connection := paginationService.BuildConnection(polls, total, hasPreviousPage, hasNextPage,
		func(poll model.Poll) uint64 {
			return poll.ID
		},
	)

	return []*model.PollConnection{connection}, nil
}

// CreatePoll creates a new poll.
func (q PollService) CreatePoll(form model.PollCreate) (*model.Poll, error) {
	// check session exists
	_, err := NewSessionService().GetSession(form.SessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	return repository.NewPollRepository().CreatePoll(form)
}

// UpdatePoll updates an existing poll.
func (q PollService) UpdatePoll(id uint64, form model.PollUpdate) (*model.Poll, error) {
	return repository.NewPollRepository().UpdatePoll(id, form)
}

// DeletePoll deletes polls based on the provided poll IDs and user information.
func (q PollService) DeletePoll(ids []uint64, userInfo model.UserInfo) ([]*model.Poll, error) {
	return repository.NewPollRepository().DeletePolls(ids, userInfo)
}
