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

type QuestionController struct {
}

func NewQuestionController() QuestionController {
	return QuestionController{}
}

// SelectOneQuestion 根據提供的 ID 檢查問題是否存在。
// @Summary
// @tags 問題
// @Summary 根據提供的 ID 檢查問題是否存在
// @Description 根據提供的 ID 檢查問題是否存在
// @Accept json
// @Produce json
// @Param id path int true "問題ID"
// @Param candidates query bool false "是否檢索候選人"
// @Success 200 {string} string "ok"
// @Router /question/{id} [get]
func (q QuestionController) SelectOneQuestion(c *gin.Context) {
	id := c.Param("id")
	questionID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid question ID",
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

	// 檢索問題及其候選人。
	var questionOne *model.Question
	candidates := c.Query("candidates")
	questionService := service.NewQuestionService()
	if candidates == "true" {
		questionOne, err = questionService.SelectQuestionWithCandidates(questionID, isAdmin, userId)
	} else {
		questionOne, err = questionService.SelectOneQuestion(questionID, isAdmin, userId)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg":    "Question not found: " + err.Error(),
			"data":   nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "Successfully retrieved question data",
			"question":   &questionOne,
		})
	}
}

// SelectAllQuestions 檢索所有問題。
// @Summary
// @tags 問題
// @Summary 檢索所有問題
// @Description 檢索所有問題
// @Accept json
// @Produce json
// @Param vote_id query int false "投票ID"
// @Success 200 {string} string "ok"
// @Router /question [get]
func (q QuestionController) SelectAllQuestions(c *gin.Context) {
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
	questions, err := service.NewQuestionService().SelectAllQuestions(voteUuid, isAdmin, userId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg":    "Questions not found: " + err.Error(),
			"data":   nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "Successfully retrieved questions data",
			"questions":   &questions,
		})
	}
}

// CreateQuestion @Summary
// @tags 問題
// @Summary 創建一個新的問題
// @Description 創建一個新的問題
// @Accept json
// @Produce json
// @Param vote_id query int true "投票ID"
// @Param title query string true "問題標題"
// @Param description query string true "問題描述"
// @Success 200 {string} string "ok"
// @Router /question [post]
func (q QuestionController) CreateQuestion(c *gin.Context) {
	var form model.QuestionCreate
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid params: " + err.Error(),
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

	if !isAdmin {
		vote, err := service.NewVoteService().SelectOneVote(form.VoteID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": -1,
				"msg":    "Vote not found: " + err.Error(),
				"data":   nil,
			})
			return
		}

		if vote.UserID != userId {
			c.JSON(http.StatusForbidden, gin.H{
				"status": -1,
				"msg":    "Permission denied",
				"data":   nil,
			})
			return
		}
	}

	question, err := service.NewQuestionService().CreateQuestion(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to create question: " + err.Error(),
			"data":   nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "Successfully create question",
			"data":   &question,
		})
	}
}

// TODO: UpdateQuestion
// UpdateQuestion @Summary
// @tags 問題
// @Summary 更新問題
// @Description 更新問題
// @Accept json
// @Produce json
// @Param vote_id query int true "投票ID"
// @Param title query string true "問題標題"
// @Param description query string true "問題描述"
// @Success 200 {string} string "ok"
// @Router /question [put]
func (q QuestionController) UpdateQuestion(c *gin.Context) {
}

// TODO: DeleteQuestion
// DeleteQuestion @Summary
// @tags 問題
// @Summary 刪除問題
// @Description 刪除問題
// @Accept json
// @Produce json
// @Param vote_id query int true "投票ID"
// @Param title query string true "問題標題"
// @Param description query string true "問題描述"
// @Success 200 {string} string "ok"
// @Router /question [delete]
func (q QuestionController) DeleteQuestion(c *gin.Context) {
}