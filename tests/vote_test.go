package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"vote/app/controller"
	// "vote/app/model"
	// "vote/app/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateVote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	voteController := controller.NewVoteController()
	router.PUT("/vote/update/:id", voteController.UpdateVote)

	t.Run("Invalid UUID format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/vote/update/invalid-uuid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Invalid UUID format: invalid UUID length: 12", response["msg"])
	})

	t.Run("Invalid params", func(t *testing.T) {
		voteId := uuid.New()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/vote/update/"+voteId.String(), bytes.NewBuffer([]byte(`{"invalid":"data"}`)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["msg"], "Invalid params")
	})

	t.Run("Vote not found", func(t *testing.T) {
		voteId := uuid.New()
		// service.NewVoteService = func() service.VoteService {
		// 	return &service.MockVoteService{
		// 		SelectOneVoteFunc: func(uuid.UUID) (model.Vote, error) {
		// 			return model.Vote{}, gorm.ErrRecordNotFound
		// 		},
		// 	}
		// }
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/vote/update/"+voteId.String(), bytes.NewBuffer([]byte(`{"title":"New Title"}`)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Vote not found "+gorm.ErrRecordNotFound.Error(), response["msg"])
	})

	t.Run("Successfully update vote", func(t *testing.T) {
		voteId := uuid.New()
		// service.NewVoteService = func() service.VoteService {
		// 	return &service.MockVoteService{
		// 		SelectOneVoteFunc: func(uuid.UUID) (model.Vote, error) {
		// 			return model.Vote{ID: voteId, UserID: 1}, nil
		// 		},
		// 		UpdateVoteFunc: func(uuid.UUID, model.VoteUpdate) (model.Vote, error) {
		// 			return model.Vote{ID: voteId, Title: "Updated Title"}, nil
		// 		},
		// 	}
		// }
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/vote/update/"+voteId.String(), bytes.NewBuffer([]byte(`{"title":"Updated Title"}`)))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Successfully update vote", response["msg"])
	})
}