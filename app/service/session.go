package service

import (
	"vote/app/model"
	"vote/app/repository"

	"github.com/google/uuid"
)

type SessionService struct {
}

func NewSessionService() SessionService {
	return SessionService{}
}

// GetSession retrieves a session by the provided UUID.
func (v SessionService) GetSession(uuid uuid.UUID) (*model.Session, error) {
	session, err := repository.NewSessionRepository().GetSessionByUUID(uuid)

	if err != nil {
		return nil, err
	} else {
		return session, nil
	}
}

// GetSessions retrieves all sessions.
func (v SessionService) GetSessions(userInfo model.UserInfo, sessionQuery *model.SessionQuery) ([]*model.SessionConnection, error) {
	sessions, total, err := repository.NewSessionRepository().GetSessions(userInfo, sessionQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.SessionQuery, model.Session]()
	sessions, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(sessions, sessionQuery)

	paginationService := NewPaginationService[model.Session, model.SessionEdge, *model.SessionConnection]()
	connection := paginationService.BuildConnection(sessions, total, hasPreviousPage, hasNextPage,
		func(session model.Session) uint64 {
			return session.ID
		},
	)

	return []*model.SessionConnection{connection}, nil
}

// CreateOneSession creates a new session.
func (v SessionService) CreateSession(form model.SessionCreate) (*model.Session, error) {
	return repository.NewSessionRepository().CreateSession(form)
}

// UpdateOneSession updates a session by the provided UUID and form data.
func (v SessionService) UpdateSession(uuid uuid.UUID, form model.SessionUpdate) (*model.Session, error) {
	return repository.NewSessionRepository().UpdateSession(uuid, form)
}

// DeleteOneSession deletes sessions by the provided UUIDs.
func (v SessionService) DeleteSession(sessionUuids []uuid.UUID) ([]*model.Session, error) {
	return repository.NewSessionRepository().DeleteSessions(sessionUuids)
}
