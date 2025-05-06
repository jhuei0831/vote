package controller

import (
	"net/http"
	"strconv"
	"vote/app/database"
	"vote/app/middleware"
	"vote/app/model"
	"vote/app/service"
	"vote/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/i18n/gi18n"
)

type UsersController struct {}

func NewUserController() UsersController {
	return UsersController{}
}

// CreateUser @Summary
// @Tags user
// @version 1.0
// @produce application/json
// @param language header string true "language"
// @param register body UserCreate true "register"
// @Success 200 string successful return value
// @Router /v1/user/create [post]
func (u UsersController) CreateUser(c *gin.Context) {
	t := gi18n.New()
	var form model.UserCreate
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
func (u UsersController) GetUser(c *gin.Context) {
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

// Login @Summary
// @Tags user
// @version 1.0
// @produce application/json
// @param register body UserLogin true "login"
// @Success 200 string successful return token
// @Router /v1/user/login [post]
func (u UsersController) Login(c *gin.Context) {
	var form model.UserLogin
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
	tokenString, refreshTokenString, err := middleware.GenUserToken(userOne.ID, form.Account, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "Failed to generate token: " + err.Error(),
			"data":   nil,
		})
		return
	}
	// Set the token in the cookie
	c.SetCookie("user-token", tokenString, 3600, "/", "", true, true)
	c.SetCookie("user-refresh_token", refreshTokenString, 3600*24*7, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Success",
		"data": gin.H{"token": tokenString, "refresh_token": refreshTokenString},
	})
}

// Logout @Summary
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @Success 200 string successful return value
// @Router /v1/user/logout [post]
func (u UsersController) Logout(c *gin.Context) {
	// Clear the token from the cookie
	c.SetCookie("user-token", "", -1, "/", "", true, true)
	c.SetCookie("user-refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Logout successful",
	})
}

// CheckAlive @Summary
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @Success 200 string successful return value
// @Router /v1/user/check-auth [post]
func (u UsersController) CheckAuth(c *gin.Context) {
	// 核心思維是前端會問後端token是否有效，如果無效就會重新登入，有效則使用refresh token取得新的token
	// 這邊的token是從cookie中取得的
	token, err := c.Cookie("user-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Authorization token not found in Cookie",
		})
		return
	}
	claims, err := middleware.ParseUserToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Invalid Token.",
		})
		return
	}
	// renew token
	tokenString, _, _ := middleware.GenUserToken(claims.ID, claims.Account, claims.Roles)
	c.SetCookie("user-token", tokenString, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Success",
		"data": gin.H{
			"id":      claims.ID,
			"account": claims.Account,
			"roles":   claims.Roles,
		},
	})		
}

// RefreshToken @Summary
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @Success 200 string successful return value
// @Router /v1/user/refresh-token [post]
func (u UsersController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Refresh token not found in Cookie",
		})
		return
	}

	// Validate the refresh token
	err = middleware.ParseRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "Invalid or expired refresh token",
		})
		return
	}

	// Generate new access token
	id := c.GetUint64("id")            // Retrieve user ID from context or database
	account := c.GetString("account")  // Retrieve account info
	roles := c.GetStringSlice("roles") // Retrieve roles

	accessToken, newRefreshToken, err := middleware.GenUserToken(id, account, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "Failed to generate new tokens",
		})
		return
	}
	// Set the token in the cookie
	c.SetCookie("user-token", accessToken, 3600, "/", "", true, true)
	c.SetCookie("user-refresh_token", newRefreshToken, 3600*24*7, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Success",
		"data": gin.H{"token": accessToken, "refresh_token": newRefreshToken},
	})
}
