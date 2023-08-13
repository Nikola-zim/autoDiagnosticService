package v1

import (
	"autoDiagnosticService/internal/entity"
	"autoDiagnosticService/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

const (
	userKey = "user"
)

type AuthHandlers struct {
	UseCase Recognition
	l       logger.Interface
}

func NewAuthHandlers(handler *gin.RouterGroup, useCase Recognition, l logger.Interface) {
	au := &AuthHandlers{
		l:       l,
		UseCase: useCase,
	}
	// LoginWEB and logout routes
	u := handler.Group("")
	{
		u.GET("/main_page", au.authPage)
		u.POST("/register", au.Register)
		u.POST("/login", au.LoginWEB)
		u.GET("/logout", au.logout)
	}
}

func (au *AuthHandlers) authPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"block_title": "Test page",
		"URL":         "/v1/auth/login",
	})
}

func (au *AuthHandlers) Register(c *gin.Context) {
	var user entity.User
	user.Login = c.PostForm("username")
	user.Password = c.PostForm("password")

	// Validate form input
	if strings.Trim(user.Login, " ") == "" || strings.Trim(user.Password, " ") == "" {
		errorResponse(c, http.StatusBadRequest, "Parameters can't be empty")
		return
	}
	if !isValidUsername(user.Login) {
		errorResponse(c, http.StatusBadRequest, "Invalid username format")
		return
	}

	err := au.UseCase.AddUser(c, user)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Registration failed")
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered user"})
}

func (au *AuthHandlers) LoginWEB(c *gin.Context) {

	var user entity.User
	user.Login = c.PostForm("username")
	user.Password = c.PostForm("password")

	// Validate form input
	if strings.Trim(user.Login, " ") == "" || strings.Trim(user.Password, " ") == "" {
		errorResponse(c, http.StatusBadRequest, "Parameters can't be empty")
		return
	}

	ok, err := au.UseCase.Login(c, user)
	if err != nil || !ok {
		errorResponse(c, http.StatusUnauthorized, "Authentication failed")
	} else {
		//Устанавливаем куку
		session := sessions.Default(c)
		session.Set(userKey, user.Login)

		if err = session.Save(); err != nil {
			errorResponse(c, http.StatusInternalServerError, "Failed to save session")
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
	}
}

func (au *AuthHandlers) logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		errorResponse(c, http.StatusBadRequest, "Invalid session token")
		return
	}
	session.Delete(userKey)
	if err := session.Save(); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Failed to save session")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func isValidUsername(username string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_]{4,16}$`)
	return regex.MatchString(username)
}
