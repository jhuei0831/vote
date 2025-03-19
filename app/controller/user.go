package controller

import (
	"net/http"
	"strconv"
	"vote/app/database"
	"vote/app/middleware"
	"vote/app/service"
	"vote/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/i18n/gi18n"
)

type UsersController struct {}

func NewUserController() UsersController {
	return UsersController{}
}

type UserRegister struct {
	Account string `json:"account" binding:"required" example:"account"`
	Password string `json:"password" binding:"required" example:"password"`
	Email string `json:"email" binding:"required" example:"test123@gmail.com"`
}

type UserLogin struct {
	Account string `json:"account" binding:"required" example:"account"`
	Password string `json:"password" binding:"required" example:"password"`
}

// CreateUser @Summary
// @Tags user
// @version 1.0
// @produce application/json
// @param language header string true "language"
// @param register body UserRegister true "register"
// @Success 200 string successful return value
// @Router /v1/user/create [post]
func (u UsersController) CreateUser (c *gin.Context) {
	t := gi18n.New()
	var form UserRegister
	bindErr := c.BindJSON(&form)

	lan := c.Request.Header.Get("language")
	if lan == "" {
		lan = "en"
	}
	t.SetLanguage(lan)
	
	if bindErr == nil {
		err := service.NewUserService().RegisterOneUser(form.Account, form.Password, form.Email)
		if err == nil {
			// go service.NewSmtpService().MultiSend(form.Email)
			c.JSON(http.StatusOK, gin.H{
				"status": 1,
				"msg": t.Translate(c, "Response_Success"),
				"data": nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -1,
				"msg": "Register Failed: " + err.Error(),
				"data": nil,
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg": "Failed to parse register data: " + utils.ValidationErrorMessage(bindErr),
			"data": nil,
		})
	}
}

// GetUser @Summary
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param id path int true "id" default(1)
// @Success 200 string successful return data
// @Router /v1/user/{id} [get]
func (u UsersController) GetUser (c *gin.Context) {
	id := c.Params.ByName("id")

	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg": "Failed to parse params" + err.Error(),
			"data": nil,
		})
	}
	userOne, err := service.NewUserService().SelectOneUsers(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg": "User not found" + err.Error(),
			"data": nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":  "Successfully get user data",
			"user": &userOne,
		})
	}
}

// AuthHandler @Summary
// @Tags user
// @version 1.0
// @produce application/json
// @param register body UserLogin true "login"
// @Success 200 string successful return token
// @Router /v1/user/login [post]
func (u UsersController) AuthHandler(c *gin.Context) {
	var form UserLogin
	bindErr := c.BindJSON(&form)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid params",
		})
		return
	}
	userOne, err := service.NewUserService().LoginOneUser(form.Account, form.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to parse params: " + err.Error(),
			"data":   nil,
		})
		return
	}

	if userOne == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg":    "User not found",
			"data":   nil,
		})
		return
	}
	roles, err := database.Enforcer.GetRolesForUser(strconv.FormatUint(userOne.ID, 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to get roles: " + err.Error(),
			"data":   nil,
		})
		return
	}
	tokenString, _ := middleware.GenToken(userOne.ID, form.Account, roles)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Success",
		"data": gin.H{"token": tokenString},
	})
}