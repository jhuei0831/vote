package repository

import (
	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type PollRepository struct {
}

func NewPollRepository() PollRepository {
	return PollRepository{}
}

// GetPollByID checks if a poll exists based on the provided ID and preloads options as needed.
func (q PollRepository) GetPoll(id uint64, isAdmin bool, userId uint64, preloadPollOptions bool) (*model.Poll, error) {
	var poll *model.Poll

	query := database.SqlSession.Where("polls.id = ?", id)

	// If poll options need to be preloaded, add them to the query.
	if preloadPollOptions {
		query = query.Preload("PollOptions")
	}

	// If the user is not an admin, add a condition for the user ID.
	if !isAdmin {
		query = query.
			Joins("JOIN sessions ON polls.session_id = sessions.uuid").
			Where("sessions.user_id = ?", userId)
	}

	err := query.First(&poll).Error
	if err != nil {
		return nil, err
	}

	return poll, nil
}

// GetPollsBySessionID retrieves polls based on the provided SessionID.
func (q PollRepository) GetPollsBySessionID(sessionID uuid.UUID) ([]model.Poll, error) {
	var polls []model.Poll

	err := database.SqlSession.
		Where("session_id = ?", sessionID).
		Find(&polls).Error
	if err != nil {
		return nil, err
	}

	return polls, nil
}

// GetPolls retrieves a list of polls based on the provided query parameters, along with the total count of matching records.
func (q PollRepository) GetPolls(pollQuery *model.PollQuery) ([]model.Poll, int64, error) {
	var polls []model.Poll
	var total int64

	query := database.SqlSession.Model(&model.Poll{}).Where("session_id = ?", pollQuery.SessionID)

	// Title LIKE filter
	if pollQuery.Title != "" {
		query = query.Where("polls.title LIKE ?", "%"+pollQuery.Title+"%")
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Use pagination service to handle pagination
	paginationRepository := NewPaginationRepository[*model.PollQuery, model.Poll]()
	query, err = paginationRepository.Handler(query, pollQuery)
	if err != nil {
		return nil, 0, err
	}

	if pollQuery.PollOptions {
		if err := query.Preload("PollOptions").Find(&polls).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := query.Find(&polls).Error; err != nil {
			return nil, 0, err
		}
	}

	// Query data
	err = query.Find(&polls).Error

	return polls, total, err
}

// CreatePoll creates a new poll.
func (q PollRepository) CreatePoll(form model.PollCreate) (*model.Poll, error) {
	poll := model.Poll{
		SessionID:   form.SessionID,
		Title:       form.Title,
		Description: form.Description,
	}

	insertErr := database.SqlSession.Model(&model.Poll{}).Create(&poll).Error

	return &poll, insertErr
}

// UpdatePoll updates an existing poll.
func (q PollRepository) UpdatePoll(id uint64, form model.PollUpdate) (*model.Poll, error) {
	var poll model.Poll

	updateError := database.SqlSession.
		Model(&poll).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Omit("session_id").
		Updates(&form).Error

	return &poll, updateError
}

// DeletePolls deletes polls.
func (q PollRepository) DeletePolls(ids []uint64, userInfo model.UserInfo) ([]*model.Poll, error) {
	var polls []*model.Poll

	query := database.SqlSession.
		Model(&model.Poll{}).Where("id IN ?", ids)

	// If the user is not an admin, check the user's ownership
	if !userInfo.IsAdmin {
		query = query.
			Where("session_id IN (SELECT uuid FROM sessions WHERE user_id = ?)", userInfo.UserID)
	}

	// Use Returning clause to delete and return the deleted records
	deleteErr := query.Clauses(clause.Returning{}).Delete(&polls).Error

	return polls, deleteErr
}
