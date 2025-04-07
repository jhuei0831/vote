package controller

import (
	"net/http"
	"strconv"
	"vote/app/database"
	"vote/app/model"
	"vote/app/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CandidateController 候選人控制器。
type CandidateController struct {
}

// NewCandidateController 創建新的候選人控制器。
func NewCandidateController() CandidateController {
	return CandidateController{}
}

// SelectOneCandidate 根據提供的 ID 檢查候選人是否存在。
// @Summary
// @tags 候選人
// @Summary 根據提供的 ID 檢查候選人是否存在
// @Description 根據提供的 ID 檢查候選人是否存在
// @Accept json
// @Produce json
// @Param id path int true "候選人ID"
// @Success 200 {string} string "ok"
// @Router /question/{question_id}/candidate/{id} [get]
func (ca CandidateController) SelectOneCandidate(c *gin.Context) {
	id := c.Param("id")
	candidateID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid candidate ID",
			"data":   nil,
		})
		return
	}

	userId := c.MustGet("id").(uint64)
	isAdmin, err := database.CheckIfAdmin(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to check user role: " + err.Error(),
			"data":   nil,
		})
		return
	}

	candidateService := service.NewCandidateService()
	candidateOne, err := candidateService.SelectOneCandidate(candidateID, isAdmin, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to select candidate: " + err.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   candidateOne,
	})
}

// SelectAllCandidates 檢索所有候選人。
// @Summary
// @tags 候選人
// @Summary 檢索所有候選人
// @Description 檢索所有候選人
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /question/{question_id}/candidate [get]
func (ca CandidateController) SelectAllCandidates(c *gin.Context) {
	voteId := c.Param("vote_id")
	voteUuid, err := uuid.Parse(voteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid vote ID",
			"data":   nil,
		})
		return
	}

	userId := c.MustGet("id").(uint64)
	isAdmin, err := database.CheckIfAdmin(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to check user role: " + err.Error(),
			"data":   nil,
		})
		return
	}

	candidateService := service.NewCandidateService()
	candidates, err := candidateService.SelectAllCandidates(voteUuid, isAdmin, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to select candidates: " + err.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   candidates,
	})
}

// CreateCandidate 創建新的候選人。
// @Summary
// @tags 候選人
// @Summary 創建新的候選人
// @Description 創建新的候選人
// @Accept json
// @Produce json
// @Param id path int true "候選人ID"
// @Success 200 {string} string "ok"
// @Router /question/{question_id}/candidate [post]
func (ca CandidateController) CreateCandidate(c *gin.Context) {
	var form model.CandidateCreate
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid request: " + err.Error(),
			"data":   nil,
		})
		return
	}
	userId := c.MustGet("id").(uint64)
	isAdmin, err := database.CheckIfAdmin(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to check user role: " + err.Error(),
			"data":   nil,
		})
		return
	}

	// check if question exists
	questionId := form.QuestionID
	questionService := service.NewQuestionService()
	_, err = questionService.SelectOneQuestion(questionId, isAdmin, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to select question: " + err.Error(),
			"data":   nil,
		})
		return
	}

	candidateService := service.NewCandidateService()
	candidate, err := candidateService.CreateCandidate(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to create candidate: " + err.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   candidate,
	})
}