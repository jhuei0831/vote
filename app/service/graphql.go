package service

import (
	"context"
	"fmt"
	"vote/app/database"
	"vote/app/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type GraphqlService struct {
}

func NewGraphqlService() GraphqlService {
	return GraphqlService{}
}

func (g GraphqlService) GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}

	return gc, nil
}

// BindQuery binds query parameters from the Gin context to the provided input struct.
func (g GraphqlService) BindQuery(ctx context.Context, input interface{}) error {
	gc, err := g.GinContextFromContext(ctx)
	if err != nil {
		return err
	}

	return gc.ShouldBindQuery(input)
}

// Get UserId from Gin context
func (g GraphqlService) GetUserIdFromContext(ctx context.Context) (model.UserInfo, error) {
	userInfo := model.UserInfo{UserID: 0, IsAdmin: false}
	gc, err := g.GinContextFromContext(ctx)
	if err != nil {
		return userInfo, err
	}

	userId, exists := gc.Get("userId")
	if !exists {
		return userInfo, gqlerror.Errorf("user not exists")
	}

	userInfo.UserID = userId.(uint64)
	return userInfo, nil
}

// Get UserId and IsAdmin from Gin context
func (g GraphqlService) GetUserInfoFromContext(ctx context.Context) (model.UserInfo, error) {
	userInfo := model.UserInfo{UserID: 0, IsAdmin: false}
	gc, err := g.GinContextFromContext(ctx)
	if err != nil {
		return userInfo, err
	}

	userId, exists := gc.Get("userId")
	if !exists {
		return userInfo, gqlerror.Errorf("user not exists")
	}

	isAdmin, err := database.CheckIfAdmin(userId.(uint64))
	if err != nil {
		return userInfo, gqlerror.Errorf("failed to check user role")
	}
	
	userInfo.UserID = userId.(uint64)
	userInfo.IsAdmin = isAdmin
	return userInfo, nil
}

// Get VoterInfo from Gin context
func (g GraphqlService) GetVoterInfoFromContext(ctx context.Context) (model.VoterInfo, error) {
	voterInfo := model.VoterInfo{VoterID: 0, SessionID: uuid.UUID{}, IsVoted: false}
	gc, err := g.GinContextFromContext(ctx)
	if err != nil {
		return voterInfo, err
	}

	voterId, exists := gc.Get("voterId")
	if !exists {
		return voterInfo, gqlerror.Errorf("voter not exists")
	}

	voteId, exists := gc.Get("voterVoteId")
	if !exists {
		return voterInfo, gqlerror.Errorf("vote not exists")
	}

	isVoted, exists := gc.Get("voterIsVoted")
	if !exists {
		return voterInfo, gqlerror.Errorf("voting status not exists")
	}

	voterInfo.VoterID = voterId.(uint64)
	voterInfo.SessionID = voteId.(uuid.UUID)
	voterInfo.IsVoted = isVoted.(bool)
	
	return voterInfo, nil
}