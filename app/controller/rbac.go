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
	creator := string(enum.Creator)
	actions := []string{"create", "read", "update", "delete"}
	resources := []string{"vote", "question", "candidate", "password", "ballot"}

	for _, res := range resources {
		for _, act := range actions {
			if _, err = enforcer.AddPolicy(creator, res, act); err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}

	// Admin role
	admin := string(enum.Admin)
	for _, act := range actions {
		if _, err = enforcer.AddPolicy(admin, "user", act); err != nil {
			break
		}
	}
	enforcer.AddRoleForUser(admin, creator)
	
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