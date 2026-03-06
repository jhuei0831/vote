package service

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/repository"
)

type PollOptionService struct {
}

func NewPollOptionService() PollOptionService {
	return PollOptionService{}
}

// GetPollOption retrieves a poll option by the provided ID, checking if the user is an admin or if they have access to the poll option based on their user ID.
func (c PollOptionService) GetPollOption(id uint64, isAdmin bool, userId uint64) (*model.PollOption, error) {
	pollOption := &model.PollOption{}
	query := database.SqlSession.
		Where("poll_options.id = ?", id)
	
	if !isAdmin {
		query = query.
			Joins("JOIN polls ON poll_options.poll_id = polls.id").
			Joins("JOIN sessions ON polls.session_id = sessions.uuid").
			Where("sessions.user_id = ?", userId)
	}
		
	err	:= query.First(&pollOption).Error
	if err != nil {
		return nil, err
	}
	return pollOption, nil
}

// GetPollOptions retrieves all poll options based on the provided query parameters.
func (c PollOptionService) GetPollOptions(pollOptionQuery *model.PollOptionQuery) ([]*model.PollOptionConnection, error) {
	pollOptions, total, err := repository.NewPollOptionRepository().GetPollOptions(pollOptionQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.PollOptionQuery, model.PollOption]()
	pollOptions, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(pollOptions, pollOptionQuery)

	paginationService := NewPaginationService[model.PollOption, model.PollOptionEdge, *model.PollOptionConnection]()
	connection := paginationService.BuildConnection(pollOptions, total, hasPreviousPage, hasNextPage,
		func(pollOption model.PollOption) uint64 {
			return pollOption.ID
		},
	)

	return []*model.PollOptionConnection{connection}, nil
}

// CreatePollOption creates a new poll option based on the provided user information and form data.
func (c PollOptionService) CreatePollOption(userInfo model.UserInfo, form model.PollOptionCreate) (*model.PollOption, error) {
	pollOption, err := repository.NewPollOptionRepository().CreatePollOption(userInfo, form)
	if err != nil {
		return nil, err
	}

	return pollOption, nil
}

func (c PollOptionService) UpdatePollOption(userInfo model.UserInfo, id uint64, form model.PollOptionUpdate) (*model.PollOption, error) {
	pollOption, err := repository.NewPollOptionRepository().UpdatePollOption(userInfo, id, form)
	if err != nil {
		return nil, err
	}

	return pollOption, nil
}

// DeletePollOption deletes poll options based on the provided IDs and user information.
func (c PollOptionService) DeletePollOption(ids []uint64, userInfo model.UserInfo) ([]*model.PollOption, error) {
	return repository.NewPollOptionRepository().DeletePollOptions(ids, userInfo)
}