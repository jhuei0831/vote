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
	// 產生Token
	tokenString, refreshTokenString, err := middleware.GenVoterToken(password.ID, voteUUID)
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
		"data": gin.H{"token": tokenString, "refresh_token": refreshTokenString},
	})
}