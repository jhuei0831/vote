package controller

import (
	"net/http"
	"vote/app/database"
	"vote/app/enum"

	"github.com/gin-gonic/gin"
)

type RbacController struct {
}

func NewRbacController() RbacController {
	return RbacController{}
}

// Rbac Initial
// @Summary
// @tags RBAC
// @Summary RBAC Initial
// @Description RBAC Initial
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /rbac/init [get]
func (r RbacController) Initial(c *gin.Context) {
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -1,
			"msg":    "User ID not found in context",
			"data":   nil,
		})
		return
	}
	isAdmin, err := database.CheckIfAdmin(userId.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to check user role: " + err.Error(),
			"data":   nil,
		})
		return
	}
	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -1,
			"msg":    "Permission denied",
			"data":   nil,
		})
		return
	}
	// Create admin role
	_, enforcer, err := database.Rbac()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    err,
			"data":   nil,
		})
		return
	}
	// Creator role
	creator := enum.Creator
	// vote
	enforcer.AddPolicy(creator, "vote", "create")
	enforcer.AddPolicy(creator, "vote", "read")
	enforcer.AddPolicy(creator, "vote", "update")
	enforcer.AddPolicy(creator, "vote", "delete")
	// question
	enforcer.AddPolicy(creator, "question", "create")
	enforcer.AddPolicy(creator, "question", "read")
	enforcer.AddPolicy(creator, "question", "update")
	enforcer.AddPolicy(creator, "question", "delete")
	// candidate
	enforcer.AddPolicy(creator, "candidate", "create")
	enforcer.AddPolicy(creator, "candidate", "read")
	enforcer.AddPolicy(creator, "candidate", "update")
	enforcer.AddPolicy(creator, "candidate", "delete")
	// password
	enforcer.AddPolicy(creator, "password", "create")
	enforcer.AddPolicy(creator, "password", "read")
	enforcer.AddPolicy(creator, "password", "update")
	enforcer.AddPolicy(creator, "password", "delete")
	// ballot
	enforcer.AddPolicy(creator, "ballot", "create")
	enforcer.AddPolicy(creator, "ballot", "read")
	enforcer.AddPolicy(creator, "ballot", "update")
	enforcer.AddPolicy(creator, "ballot", "delete")
	
	// Admin role
	admin := enum.Admin
	enforcer.AddPolicy(admin, "user", "create")
	enforcer.AddPolicy(admin, "user", "read")
	enforcer.AddPolicy(admin, "user", "update")
	enforcer.AddPolicy(admin, "user", "delete")
	enforcer.AddRoleForUser(string(admin), string(creator))
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    err,
			"data":   nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "Successfully init RBAC",
		"data":   nil,
	})
}