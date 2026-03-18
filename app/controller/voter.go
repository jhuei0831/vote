package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vote/app/middleware"
	"vote/app/model"
	"vote/app/repository"
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

// VoterLogin Anonymous voter login
// @Summary Anonymous voter login
// @tags Anonymous Voting
// @Summary Anonymous voter login
// @Description Anonymous voter login
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /voter/login [post]
func (a VoterController) VoterLogin(c *gin.Context) {
	// Set timeout control
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var form model.VoterVerify
	if err := c.BindJSON(&form); err != nil {
		utils.HandleError(c, http.StatusBadRequest, -1, "Invalid params", err)
		return
	}

	// Validate UUID
	voteUUID, err := uuid.Parse(form.SessionID.String())
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, -1, "Invalid session ID", err)
		return
	}

	// Define result structures for concurrent processing
	type passwordResult struct {
		invitation *model.Invitation
		err        error
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

	// Create all channels with appropriate size to avoid goroutine leaks
	passwordCh := make(chan passwordResult, 1)
	votedCh := make(chan votedResult, 1)
	tokenCh := make(chan tokenResult, 1)

	// Password encryption and validation - single goroutine handles entire flow
	go func() {
		// Encrypt password
		passwordEncrypt, err := (&utils.Invitation{}).Encrypt(form.CodeHash)
		if err != nil {
			passwordCh <- passwordResult{nil, fmt.Errorf("failed to encrypt password: %w", err)}
			return
		}

		// Check password
		password, err := service.NewInvitationService().GetInvitation(voteUUID, passwordEncrypt)
		fmt.Println(passwordEncrypt)
		if err != nil {
			passwordCh <- passwordResult{nil, fmt.Errorf("failed to validate password: %w", err)}
			return
		}

		passwordCh <- passwordResult{password, nil}
	}()

	// Receive password check result with timeout handling
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
		voter = res.invitation.ID
	}

	// Check if user has already voted
	go func() {
		hasVoted, err := repository.NewBallotRepository().CheckIfVoterHasVoted(voter)
		votedCh <- votedResult{hasVoted, err}
	}()

	// Receive voting status check result with timeout handling
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

	// Generate Token
	go func() {
		tokenString, refreshToken, err := middleware.GenVoterToken(voter, voteUUID, isVoted)
		tokenCh <- tokenResult{tokenString, refreshToken, err}
	}()

	// Receive Token generation result with timeout handling
	select {
	case <-ctx.Done():
		utils.HandleError(c, http.StatusGatewayTimeout, -1, "Request timeout during token generation", nil)
		return
	case res := <-tokenCh:
		if res.err != nil {
			utils.HandleError(c, http.StatusInternalServerError, -1, "Failed to generate authentication tokens", res.err)
			return
		}

		// Store token in cookie
		c.SetCookie("voter-token", res.token, 3600, "/", "", true, true)

		// Return success response
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Voter login success",
			"data": gin.H{
				"token": res.token,
			},
		})
	}
}

// CheckAuth Check voter's Token
// @Summary Check voter's Token
// @tags Anonymous Voting
// @Description Check voter's Token
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /voter/check-auth [post]
func (a VoterController) CheckAuth(c *gin.Context) {
	// Set timeout control
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
	}

	// Check voting status and generate new token
	type authResult struct {
		isVoted     bool
		tokenString string
		err         error
	}
	resultCh := make(chan authResult, 1)

	go func() {
		// Check voting status
		hasVoted, err := repository.NewBallotRepository().CheckIfVoterHasVoted(claims.ID)
		if err != nil {
			resultCh <- authResult{false, "", fmt.Errorf("failed to check voting status: %w", err)}
			return
		}

		// Regenerate token
		tokenString, _, err := middleware.GenVoterToken(claims.ID, claims.SessionID, hasVoted)
		if err != nil {
			resultCh <- authResult{hasVoted, "", fmt.Errorf("failed to generate token: %w", err)}
			return
		}

		resultCh <- authResult{hasVoted, tokenString, nil}
	}()

	// Process result with timeout control
	select {
	case <-ctx.Done():
		utils.HandleError(c, http.StatusGatewayTimeout, -1, "Request timeout during auth check", nil)
		return
	case res := <-resultCh:
		if res.err != nil {
			utils.HandleError(c, http.StatusInternalServerError, -1, "Authentication check failed", res.err)
			return
		}

		// Update token in cookie
		c.SetCookie("voter-token", res.tokenString, 3600, "/", "", true, true)

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{
				"id":        claims.ID,
				"sessionID": claims.SessionID,
				"voted":     res.isVoted,
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

// CheckIsVoted Check if voter has already voted
// @Summary Check if voter has already voted
// @tags Anonymous Voting
// @Description Check if voter has already voted
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /voter/is-voted [get]
func (a VoterController) CheckIsVoted(c *gin.Context) {
	// Set timeout control
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

	// Use goroutine to check voting status
	type checkResult struct {
		hasVoted bool
		err      error
	}
	resultCh := make(chan checkResult, 1)

	go func() {
		hasVoted, err := repository.NewBallotRepository().CheckIfVoterHasVoted(claims.ID)
		resultCh <- checkResult{hasVoted, err}
	}()

	// Process result with timeout control
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
