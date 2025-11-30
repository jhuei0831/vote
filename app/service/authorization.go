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
	ErrVoteIDsMissing = errors.New("no vote ids provided")
	ErrForbiddenVote  = errors.New("forbidden")
)

type AuthorizationService struct{}

func NewAuthorizationService() AuthorizationService {
	return AuthorizationService{}
}

// EnsureVoteOwnership verifies that the user either has admin privileges or owns every vote ID provided.
func (a AuthorizationService) EnsureVoteOwnership(userID uint64, voteIDs []uuid.UUID) error {
	if len(voteIDs) == 0 {
		return ErrVoteIDsMissing
	}

	isAdmin, err := database.CheckIfAdmin(userID)
	if err != nil {
		return err
	}

	if isAdmin {
		return nil
	}

	var ownedCount int64
	if err := database.SqlSession.Model(&model.Vote{}).
		Where("uuid IN ?", voteIDs).
		Where("user_id = ?", userID).
		Count(&ownedCount).Error; err != nil {
		return err
	}

	if ownedCount != int64(len(voteIDs)) {
		return ErrForbiddenVote
	}

	return nil
}

// AuthorizeVoteAccess checks if the user has permission to perform an action on a vote.
func (a AuthorizationService) AuthorizeVoteAccess(ctx context.Context, voteID interface{}, action string) error {
	userInfo, err := NewGraphqlService().GetUserInfoFromContext(ctx)
	if err != nil {
		return gqlerror.Errorf("failed to get user info from context: %v", err)
	}

	if !userInfo.IsAdmin {
		var voteIDs []uuid.UUID
		switch v := voteID.(type) {
			case uuid.UUID:
				if v != uuid.Nil {
					voteIDs = []uuid.UUID{v}
				}
			case []uuid.UUID:
				voteIDs = v
			default:
				return gqlerror.Errorf("invalid voteID type")
		}

		if len(voteIDs) > 0 {
			if err := a.EnsureVoteOwnership(userInfo.UserID, voteIDs); err != nil {
				return gqlerror.Errorf("failed to authorize %s: %v", action, err)
			}
		}
	}
	
	return nil
}
