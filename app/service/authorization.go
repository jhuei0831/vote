package service

import (
	"context"
	"errors"

	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	ErrSessionIDsMissing = errors.New("no session ids provided")
	ErrForbiddenSession  = errors.New("forbidden")
)

type AuthorizationService struct{}

func NewAuthorizationService() AuthorizationService {
	return AuthorizationService{}
}

// EnsureSessionOwnership verifies that the user either has admin privileges or owns every session ID provided.
func (a AuthorizationService) EnsureSessionOwnership(userID uint64, SessionIDs []uuid.UUID) error {
	if len(SessionIDs) == 0 {
		return ErrSessionIDsMissing
	}

	isAdmin, err := database.CheckIfAdmin(userID)
	if err != nil {
		return err
	}

	if isAdmin {
		return nil
	}

	var ownedCount int64
	if err := database.SqlSession.Model(&model.Session{}).
		Where("uuid IN ?", SessionIDs).
		Where("user_id = ?", userID).
		Count(&ownedCount).Error; err != nil {
		return err
	}

	if ownedCount != int64(len(SessionIDs)) {
		return ErrForbiddenSession
	}

	return nil
}

// AuthorizeSessionAccess checks if the user has permission to perform an action on a session.
func (a AuthorizationService) AuthorizeSessionAccess(ctx context.Context, SessionID interface{}, action string) error {
	userInfo, err := NewGraphqlService().GetUserInfoFromContext(ctx)
	if err != nil {
		return gqlerror.Errorf("failed to get user info from context: %v", err)
	}

	if !userInfo.IsAdmin {
		var SessionIDs []uuid.UUID
		switch v := SessionID.(type) {
			case uuid.UUID:
				if v != uuid.Nil {
					SessionIDs = []uuid.UUID{v}
				}
			case []uuid.UUID:
				SessionIDs = v
			default:
				return gqlerror.Errorf("invalid SessionID type")
		}

		if len(SessionIDs) > 0 {
			if err := a.EnsureSessionOwnership(userInfo.UserID, SessionIDs); err != nil {
				return gqlerror.Errorf("failed to authorize %s: %v", action, err)
			}
		}
	}
	
	return nil
}
