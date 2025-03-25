package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"vote/app/database"
	"vote/app/model"
	"vote/app/service"
	"vote/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PasswordController struct {
}

func NewPasswordController() PasswordController {
	return PasswordController{}
}

// CreatePassword 建立可以加解密的密碼
// @Summary
// @tags 密碼
// @Summary 建立可以加解密的密碼
// @Description 建立可以加解密的密碼
// @Accept json
// @Produce json
// @Param vote_id path string true "投票ID"
// @Success 200 {string} string "ok"
// @Router /password/create [post]
func (p PasswordController) CreatePassword(c *gin.Context) {
	var form model.PasswordCreate
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid request: " + utils.ValidationErrorMessage(err),
			"data":   nil,
		})
		return
	}

	voteUUID, err := uuid.Parse(form.VoteID.String())
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

	voteService := service.NewVoteService()
	vote, err := voteService.SelectOneVote(voteUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to select vote: " + err.Error(),
			"data":   nil,
		})
		return
	}

	if !isAdmin && vote.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -1,
			"msg":    "Permission denied",
			"data":   nil,
		})
		return
	}

	passwordService := service.NewPasswordService()
	err = passwordService.CreatePassword(voteUUID, form.Number, form.Length, form.Format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to create password: " + err.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "successfully create password",
		"data":   form.Number,
	})
}

// DecryptPassword 解密密碼
// @Summary
// @tags 密碼
// @Summary 解密密碼
// @Description 解密密碼
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /password/decrypt [post]
func (p PasswordController) DecryptPassword(c *gin.Context) {
	jsonData, _ := io.ReadAll(c.Request.Body)
	// json to Array
	var passwords []string
	err := json.Unmarshal(jsonData, &passwords)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid JSON format: " + err.Error(),
			"data":   nil,
		})
		return
	}
	if len(passwords) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "No password provided",
			"data":   nil,
		})
		return
	}

	decrypts := make(map[string]string)
	for _, password := range passwords {
		decrypt, err := (&utils.Password{}).Decrypt(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -1,
				"msg":    "Failed to decrypt password [" + password + "]: " + err.Error(),
				"data":   nil,
			})
			return
		}
		decrypts[password] = decrypt
	}

	

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "successfully decrypt password",
		"data":   decrypts,
	})
}