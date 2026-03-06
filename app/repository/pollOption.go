package repository

import (
	"fmt"
	"vote/app/database"
	"vote/app/model"

	"gorm.io/gorm/clause"
)

type PollOptionRepository struct {
}

func NewPollOptionRepository() PollOptionRepository {
	return PollOptionRepository{}
}

// GetPollOptionByID checks whether a poll option exists based on the provided ID.
func (c PollOptionRepository) GetPollOptionByID(id uint64, pollId uint64) (*model.PollOption, error) {
	pollOption := &model.PollOption{}
	err := database.SqlSession.
		Where("id = ? AND poll_id = ?", id, pollId).
		First(&pollOption).Error
	if err != nil {
		return nil, err
	}
	return pollOption, nil
}

// GetPollOptions retrieves a list of poll options based on the provided query parameters, along with the total count of matching records.
func (c PollOptionRepository) GetPollOptions(candidateQuery *model.PollOptionQuery) ([]model.PollOption, int64, error) {
	var candidates []model.PollOption
	var total int64

	query := database.SqlSession.Model(&model.PollOption{})

	// PollID Equal filter
	if candidateQuery.PollID != 0 {
		query = query.Where("poll_id = ?", candidateQuery.PollID)
	}

	// Name LIKE filter
	if candidateQuery.Name != "" {
		query = query.Where("name LIKE ?", "%"+candidateQuery.Name+"%")
	}

	// Calculate total count before applying pagination
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Retrieve data with pagination
	paginationRepository := NewPaginationRepository[*model.PollOptionQuery, model.Poll]()
	query, err = paginationRepository.Handler(query, candidateQuery)
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&candidates).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve data
	err = query.Find(&candidates).Error

	return candidates, total, err
}

// CreatePollOption creates a new poll option.
func (c PollOptionRepository) CreatePollOption(userInfo model.UserInfo, form model.PollOptionCreate) (*model.PollOption, error) {
	// Verify that the pollID belongs to the SessionID
	pollRepository := NewPollRepository()
	_, err := pollRepository.GetPoll(form.PollID, userInfo.IsAdmin, userInfo.UserID, false)
	if err != nil {
		return nil, fmt.Errorf("%s", "poll record not found")
	}

	pollOption := model.PollOption{
		PollID: form.PollID,
		Name:   form.Name,
	}

	insertErr := database.SqlSession.Model(&model.PollOption{}).Create(&pollOption).Error

	return &pollOption, insertErr
}

func (c PollOptionRepository) UpdatePollOption(userInfo model.UserInfo, id uint64, form model.PollOptionUpdate) (*model.PollOption, error) {
	// Verify that the pollID and poll option id belongs to the poll option
	_, err := c.GetPollOptionByID(id, form.PollID)
	if err != nil {
		return nil, fmt.Errorf("%s", "poll option record not found")
	}

	var pollOption model.PollOption

	updateError := database.SqlSession.Model(&pollOption).
		Clauses(clause.Returning{}).
		Where("id=?", id).
		Omit("poll_id"). // Cannot update poll_id to prevent moving to another poll
		Updates(&form).Error

	return &pollOption, updateError
}

// DeletePollOptions deletes poll options.
func (c PollOptionRepository) DeletePollOptions(ids []uint64, userInfo model.UserInfo) ([]*model.PollOption, error) {
	var pollOptions []*model.PollOption

	query := database.SqlSession.
		Model(&model.PollOption{}).
		Where("poll_options.id IN ?", ids)

	// Non-admin users need to check ownership
	if !userInfo.IsAdmin {
		query = query.
			Joins("JOIN polls ON poll_options.poll_id = polls.id").
			Joins("JOIN sessions ON polls.session_id = sessions.uuid").
			Where("sessions.user_id = ?", userInfo.UserID)
	}

	// Retrieve poll options and delete them directly
	if err := query.Find(&pollOptions).Error; err != nil {
		return nil, err
	}

	if len(pollOptions) == 0 {
		return pollOptions, nil
	}

	// Extract IDs and delete
	pollOptionIds := make([]uint64, len(pollOptions))
	for i, pollOption := range pollOptions {
		pollOptionIds[i] = pollOption.ID
	}

	deleteError := database.SqlSession.Delete(&model.PollOption{}, pollOptionIds).Error
	return pollOptions, deleteError
}
