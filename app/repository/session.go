package repository

import (
	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type SessionRepository struct {
}

func NewSessionRepository() SessionRepository {
	return SessionRepository{}
}

// GetSessionByUUID retrieves a session by the provided UUID.
func (v SessionRepository) GetSessionByUUID(uuid uuid.UUID) (*model.Session, error) {
	session := &model.Session{}

	err := database.SqlSession.
		Where("uuid = ?", uuid).
		Preload("Polls.PollOptions").
		First(&session).Error

	return session, err
}

// GetSessions retrieves all sessions based on the provided conditions.
func (v SessionRepository) GetSessions(userInfo model.UserInfo, sessionQuery *model.SessionQuery) ([]model.Session, int64, error) {
	var sessions []model.Session
	var total int64

	query := database.SqlSession.Model(&model.Session{}).Preload("Polls")

	if !userInfo.IsAdmin {
		query = query.Where("user_id = ?", userInfo.UserID)
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Use pagination service to handle pagination
	paginationRepository := NewPaginationRepository[*model.SessionQuery, model.Session]()
	query, err = paginationRepository.Handler(query, sessionQuery)
	if err != nil {
		return nil, 0, err
	}

	// Query data
	err = query.Find(&sessions).Error

	return sessions, total, err
}

// CreateSession creates a new session.
func (v SessionRepository) CreateSession(form model.SessionCreate) (*model.Session, error) {
	session := model.Session{
		Title:       form.Title,
		Description: form.Description,
		UserID:      form.UserID,
		StartTime:   form.StartTime,
		EndTime:     form.EndTime,
	}

	insertErr := database.SqlSession.Create(&session).Error

	return &session, insertErr
}

// UpdateSession updates an existing session.
func (v SessionRepository) UpdateSession(uuid uuid.UUID, form model.SessionUpdate) (*model.Session, error) {
	var session model.Session

	updateError := database.SqlSession.Model(&session).
		Clauses(clause.Returning{}).
		Where("uuid = ?", uuid).
		Updates(&form).Error

	return &session, updateError
}

// DeleteSessions deletes sessions.
func (v SessionRepository) DeleteSessions(sessionUuids []uuid.UUID) ([]*model.Session, error) {
	var sessions []*model.Session

	deleteErr := database.SqlSession.
		Clauses(clause.Returning{}).
		Where("uuid IN ?", sessionUuids).
		Delete(&sessions).Error

	return sessions, deleteErr
}
