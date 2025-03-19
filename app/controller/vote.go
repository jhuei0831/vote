package controller

import (
	"net/http"
	"strconv"
	"vote/app/database"
	"vote/app/enum"
	"vote/app/model"
	"vote/app/service"
	"vote/app/utils"

	"github.com/gin-gonic/gin"
)

type VoteController struct {
}

func NewVoteController() VoteController {
	return VoteController{}
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
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -1,
			"msg":    "User ID not found in context",
			"data":   nil,
		})
		return
	}
	// 檢查用戶是否是管理員
	isAdmin, err := database.Enforcer.HasRoleForUser(strconv.FormatUint(userId.(uint64), 10), enum.Roles.Admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg": "Failed to check user role" + err.Error(),
			"data": nil,
		})
	}
	// 根據用戶ID及權限檢索所有投票
	votes, err := service.NewVoteService().SelectAllVotes(isAdmin, userId.(uint64))
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
	var form model.VoteCreate
	bindErr := c.BindJSON(&form)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid params" + utils.ValidationErrorMessage(bindErr),
		})
		return
	}
	// 把創建者的ID從Header的JWT中取出來
	userId, _ := c.Get("id")
	form.UserID = userId.(uint64)
	insertErr := service.NewVoteService().CreateOneVote(form)
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

// UpdateVote @Summary
// @tags 投票
// @Summary 更新一個投票
// @Description 更新一個投票
// @Accept json
// @Produce json
// @Param id path int true "投票ID"
// @Param title query string true "投票標題"
// @Param description query string true "投票描述"
// @Param startTime query string true "開始時間"
// @Param endTime query string true "結束時間"
// @Success 200 {string} string "ok"
// @Router /vote/update/{id} [put]
// TODO: 除了ADMIN，只有投票的創建者才能更新投票
func (v VoteController) UpdateVote (c *gin.Context) {
	id := c.Params.ByName("id")
	voteId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg": "Failed to parse params" + err.Error(),
			"data": nil,
		})
	}

	var form model.VoteUpdate
	bindErr := c.BindJSON(&form)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid params: " + utils.ValidationErrorMessage(bindErr),
		})
		return
	}

	updateErr := service.NewVoteService().UpdateOneVote(voteId, form)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to parse params: " + updateErr.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Successfully update vote",
		"data":   &form,
	})
}