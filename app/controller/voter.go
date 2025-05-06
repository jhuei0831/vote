package controller

import (
	"net/http"
	"vote/app/middleware"
	"vote/app/model"
	"vote/app/service"
	"vote/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VoterController struct {
}

func NewVoterController() VoterController {
	return VoterController{}
}

// VoterLogin 匿名投票登入
// @Summary 匿名投票登入
// @tags 匿名投票
// @Summary 匿名投票登入
// @Description 匿名投票登入
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /voter/login [post]
func (a VoterController) VoterLogin(c *gin.Context) {
	var form model.VoterLogin
	bindErr := c.BindJSON(&form)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid params",
		})
		return
	}

	// 驗證UUID
	voteUUID, err := uuid.Parse(form.VoteID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid vote ID",
			"data":   nil,
		})
		return
	}
	// 密碼加密
	passwordEncrypt, err := (&utils.Password{}).Encrypt(form.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to encrypt password: " + err.Error(),
			"data":   nil,
		})
		return
	}
	// 檢查密碼
	password, err := service.NewPasswordService().SelectOnePassword(voteUUID, passwordEncrypt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to parse params: " + err.Error(),
			"data":   nil,
		})
		return
	}
	voter := password.ID
	var isVoted = false
	if hasVoted, err := service.NewBallotService().CheckIfVoterHasVoted(voter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Failed to check if voter has voted: " + err.Error(),
		})
		return
	} else {
		isVoted = hasVoted
	}
	
	// 產生Token
	tokenString, _, err := middleware.GenVoterToken(password.ID, voteUUID, isVoted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "Failed to generate new tokens",
		})
		return
	}
	// Set the token in the cookie
	c.SetCookie("voter-token", tokenString, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Voter login success",
		"data": gin.H{"token": tokenString},
	})
}

// CheckAuth 檢查投票者的Token
// @Summary 檢查投票者的Token
// @tags 匿名投票
// @Description 檢查投票者的Token
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /voter/check-auth [post]
func (a VoterController) CheckAuth(c *gin.Context) {
	token, err := c.Cookie("voter-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Authorization token not found in Cookie",
		})
		return
	}
	claims, err := middleware.ParseVoterToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Invalid Token.",
		})
		return
	}
	// 重新產生 token
	var isVoted = false
	if hasVoted, err := service.NewBallotService().CheckIfVoterHasVoted(claims.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Failed to check if voter has voted: " + err.Error(),
		})
		return
	} else {
		isVoted = hasVoted
	}
	tokenString, _, err := middleware.GenVoterToken(claims.ID, claims.VoteID, isVoted)
	c.SetCookie("voter-token", tokenString, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Success",
		"data": gin.H{
			"id":           claims.ID,
			"voteId":       claims.VoteID,
			"voted":        isVoted,
		},
	})
}

// Logout @Summary
// @Tags voter
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @Success 200 string successful return value
// @Router /v1/voter/logout [post]
func (a VoterController) Logout(c *gin.Context) {
	// Clear the token from the cookie
	c.SetCookie("voter-token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Logout successful",
	})
}

// CheckIsVoted 檢查投票者是否已經投票
// @Summary 檢查投票者是否已經投票
// @tags 匿名投票
// @Description 檢查投票者是否已經投票
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /voter/is-voted [get]
func (a VoterController) CheckIsVoted(c *gin.Context) {
	token, err := c.Cookie("voter-token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Authorization token not found in Cookie",
		})
		return
	}
	claims, err := middleware.ParseVoterToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Invalid Token.",
		})
		return
	}

	ballotService := service.NewBallotService()
	voter := claims.ID
	if hasVoted, err := ballotService.CheckIfVoterHasVoted(voter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Failed to check if voter has voted: " + err.Error(),
		})
		return
	} else if hasVoted {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Voter has already voted.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Voter has not voted yet",
	})
}
