package service

import (
	"context"
	"fmt"
	"vote/app/database"

	"github.com/gin-gonic/gin"
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

	return gc.BindQuery(input)
}

// Get UserId from Gin context
func (g GraphqlService) GetUserIdFromContext(ctx context.Context) (uint64, error) {
	gc, err := g.GinContextFromContext(ctx)
	if err != nil {
		return 0, err
	}

	userId, exists := gc.Get("id")
	if !exists {
		return 0, gqlerror.Errorf("user not exists")
	}

	return userId.(uint64), nil
}

// Get UserId and IsAdmin from Gin context
func (g GraphqlService) GetUserInfoFromContext(ctx context.Context) (uint64, bool, error) {
	gc, err := g.GinContextFromContext(ctx)
	if err != nil {
		return 0, false, err
	}

	userId, exists := gc.Get("id")
	if !exists {
		return 0, false, gqlerror.Errorf("user not exists")
	}

	isAdmin, err := database.CheckIfAdmin(userId.(uint64))
	if err != nil {
		return 0, false, gqlerror.Errorf("failed to check user role")
	}

	return userId.(uint64), isAdmin, nil
}