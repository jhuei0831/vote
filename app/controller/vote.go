package controller

import (
	"net/http"
	"strconv"
	"vote/app/service"

	"github.com/gin-gonic/gin"
)

type VoteController struct {
}

func NewVoteController() VoteController {
	return VoteController{}
}

type VoteCreate struct {
	Title       string `json:"title" binding:"required" example:"title"`
	Description string `json:"description" binding:"required" example:"description"`
	StartTime   string `json:"startTime" binding:"required" example:"2021-01-01 00:00:00"`
	EndTime     string `json:"endTime" binding:"required" example:"2021-01-01 00:00:00"`
}

// SelectOneVote 根據提供的 ID 檢查投票是否存在。
// @Summary
// @tags 投票
// @Summary 根據提供的 ID 檢查投票是否存在
// @Description 根據提供的 ID 檢查投票是否存在
// @Accept json
// @Produce json
// @Param id path int true "投票ID"
// @Success 200 {string} string "ok"
// @Router /vote/{id} [get]
func (v VoteController) SelectOneVote (c *gin.Context) {
	id := c.Params.ByName("id")
	voteId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg": "Failed to parse params" + err.Error(),
			"data": nil,
		})
	}

	voteOne, err := service.NewVoteService().SelectOneVote(voteId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg": "Vote not found " + err.Error(),
			"data": nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":  "Successfully get vote data",
			"vote": &voteOne,
		})
	}
}

// SelectAllVotes 檢索所有投票。
// @Summary
// @tags 投票
// @Summary 檢索所有投票
// @Description 檢索所有投票
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /vote/all [get]
func (v VoteController) SelectAllVotes (c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	userId, _ := c.Get("id")
	votes, err := service.NewVoteService().SelectAllVotes(isAdmin.(bool), userId.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg": "Vote not found " + err.Error(),
			"data": nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":  "Successfully get vote data",
			"vote": &votes,
		})
	}
}

// CreateVote @Summary
// @tags 投票
// @Summary 創建一個新的投票
// @Description 創建一個新的投票
// @Accept json
// @Produce json
// @Param title query string true "投票標題"
// @Param description query string true "投票描述"
// @Param startTime query string true "開始時間"
// @Param endTime query string true "結束時間"
// @Success 200 {string} string "ok"
// @Router /vote/create [post]
func (v VoteController) CreateVote (c *gin.Context) {
	var form VoteCreate
	bindErr := c.BindJSON(&form)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid params",
		})
		return
	}
	// 把創建者的ID從Header的JWT中取出來
	creatorId, _ := c.Get("id")

	insertErr := service.NewVoteService().
		CreateOneVote(form.Title, form.Description, creatorId.(uint64), form.StartTime, form.EndTime)
	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to parse params: " + insertErr.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Successfully create vote",
		"data":   &form,
	})
}