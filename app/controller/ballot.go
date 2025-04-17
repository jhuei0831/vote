package controller

import (
	"vote/app/enum"
	"vote/app/middleware"
	"vote/app/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

type BallotController struct {
}

func NewBallotController() BallotController {
	return BallotController{}
}

// CreateBallots 建立投票
// @Summary
// @tags 投票
// @Summary 建立投票
// @Description 建立投票
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /ballot/create [post]
func (b BallotController) CreateBallots(c *gin.Context) {
	var ballots map[uint64]map[uint64]bool
	if err := c.ShouldBindJSON(&ballots); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Invalid JSON format: " + err.Error(),
			"data":   nil,
		})
		return
	}

	token, err := c.Cookie("voter-token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": enum.VoterNotLoggedIn,
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

	if err := ballotService.CreateBallots(voter, ballots); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to create ballots: " + err.Error(),
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"msg":     "Vote successfully",
		"voter":   claims.ID,
		"voteId":  claims.VoteID,
	})
}
