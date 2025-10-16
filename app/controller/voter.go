package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"
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
	// 設定超時控制
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var form model.VoterLogin
	if err := c.BindJSON(&form); err != nil {
		utils.HandleError(c, http.StatusBadRequest, -1, "Invalid params", err)
		return
	}

	// 驗證UUID
	voteUUID, err := uuid.Parse(form.VoteID.String())
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, -1, "Invalid vote ID", err)
		return
	}

	// 定義用於並行處理的結果結構
	type passwordResult struct {
		password *model.Password
		err      error
	}

	type votedResult struct {
		isVoted bool
		err     error
	}

	type tokenResult struct {
		token   string
		refresh string
		err     error
	}

	// 創建所有channel並設定適當大小以避免goroutine洩漏
	passwordCh := make(chan passwordResult, 1)
	votedCh := make(chan votedResult, 1)
	tokenCh := make(chan tokenResult, 1)

	// 密碼加密與驗證 - 單一goroutine處理整個流程
	go func() {
		// 密碼加密
		passwordEncrypt, err := (&utils.Password{}).Encrypt(form.Password)
		if err != nil {
			passwordCh <- passwordResult{nil, fmt.Errorf("failed to encrypt password: %w", err)}
			return
		}

		// 密碼檢查
		password, err := service.NewPasswordService().SelectOnePassword(voteUUID, passwordEncrypt)
		if err != nil {
			passwordCh <- passwordResult{nil, fmt.Errorf("failed to validate password: %w", err)}
			return
		}

		passwordCh <- passwordResult{password, nil}
	}()

	// 接收密碼檢查結果，附帶超時處理
	var voter uint64
	select {
	case <-ctx.Done():
		utils.HandleError(c, http.StatusGatewayTimeout, -1, "Request timeout during password validation", nil)
		return
	case res := <-passwordCh:
		if res.err != nil {
			utils.HandleError(c, http.StatusBadRequest, -1, "Authentication failed", res.err)
			return
		}
		voter = res.password.ID
	}

	// 檢查用戶是否已經投票
	go func() {
		hasVoted, err := service.NewBallotService().CheckIfVoterHasVoted(voter)
		votedCh <- votedResult{hasVoted, err}
	}()

	// 接收投票狀態檢查結果，附帶超時處理
	var isVoted bool
	select {
	case <-ctx.Done():
		utils.HandleError(c, http.StatusGatewayTimeout, -1, "Request timeout during voting status check", nil)
		return
	case res := <-votedCh:
		if res.err != nil {
			utils.HandleError(c, http.StatusBadRequest, -1, "Failed to check voting status", res.err)
			return
		}
		isVoted = res.isVoted
	}

	if isVoted {
		utils.HandleError(c, http.StatusBadRequest, -1, "Voter has already voted", nil)
		return
	}

	// 產生Token
	go func() {
		tokenString, refreshToken, err := middleware.GenVoterToken(voter, voteUUID, isVoted)
		tokenCh <- tokenResult{tokenString, refreshToken, err}
	}()

	// 接收Token生成結果，附帶超時處理
	select {
	case <-ctx.Done():
		utils.HandleError(c, http.StatusGatewayTimeout, -1, "Request timeout during token generation", nil)
		return
	case res := <-tokenCh:
		if res.err != nil {
			utils.HandleError(c, http.StatusInternalServerError, -1, "Failed to generate authentication tokens", res.err)
			return
		}

		// 將token存入cookie
		c.SetCookie("voter-token", res.token, 3600, "/", "", true, true)

		// 返回成功響應
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Voter login success",
			"data": gin.H{
				"token": res.token,
			},
		})
	}
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
	// 設定超時控制
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	token, err := c.Cookie("voter-token")
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, -1, "Authorization token not found in Cookie", nil)
		return
	}

	claims, err := middleware.ParseVoterToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, -1, "Invalid Token", err)
		return
	} // 檢查投票狀態並生成新token
	type authResult struct {
		isVoted     bool
		tokenString string
		err         error
	}
	resultCh := make(chan authResult, 1)

	go func() {
		// 檢查投票狀態
		hasVoted, err := service.NewBallotService().CheckIfVoterHasVoted(claims.ID)
		if err != nil {
			resultCh <- authResult{false, "", fmt.Errorf("failed to check voting status: %w", err)}
			return
		}

		// 重新產生token
		tokenString, _, err := middleware.GenVoterToken(claims.ID, claims.VoteID, hasVoted)
		if err != nil {
			resultCh <- authResult{hasVoted, "", fmt.Errorf("failed to generate token: %w", err)}
			return
		}

		resultCh <- authResult{hasVoted, tokenString, nil}
	}()

	// 處理結果並添加超時控制
	select {
	case <-ctx.Done():
		utils.HandleError(c, http.StatusGatewayTimeout, -1, "Request timeout during auth check", nil)
		return
	case res := <-resultCh:
		if res.err != nil {
			utils.HandleError(c, http.StatusInternalServerError, -1, "Authentication check failed", res.err)
			return
		}

		// 更新cookie中的token
		c.SetCookie("voter-token", res.tokenString, 3600, "/", "", true, true)

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{
				"id":     claims.ID,
				"voteId": claims.VoteID,
				"voted":  res.isVoted,
			},
		})
	}
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
	// 設定超時控制
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	token, err := c.Cookie("voter-token")
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, -1, "Authorization token not found in Cookie", nil)
		return
	}

	claims, err := middleware.ParseVoterToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, -1, "Invalid Token", err)
		return
	}

	// 使用goroutine檢查投票狀態
	type checkResult struct {
		hasVoted bool
		err      error
	}
	resultCh := make(chan checkResult, 1)

	go func() {
		hasVoted, err := service.NewBallotService().CheckIfVoterHasVoted(claims.ID)
		resultCh <- checkResult{hasVoted, err}
	}()

	// 處理結果並添加超時控制
	select {
	case <-ctx.Done():
		utils.HandleError(c, http.StatusGatewayTimeout, -1, "Request timeout during voting status check", nil)
		return
	case res := <-resultCh:
		if res.err != nil {
			utils.HandleError(c, http.StatusBadRequest, -1, "Failed to check voting status", res.err)
			return
		}

		if res.hasVoted {
			utils.HandleError(c, http.StatusBadRequest, -1, "Voter has already voted", nil)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Voter has not voted yet",
		})
	}
}
