package v1

import (
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	userKey = "user"
)

type AuthHandlers struct {
	useCase usecase.Recognition
	l       logger.Interface
}

func NewAuthHandlers(handler *gin.RouterGroup, useCase usecase.Recognition, l logger.Interface) {
	au := &AuthHandlers{
		l:       l,
		useCase: useCase,
	}
	// Login and logout routes
	u := handler.Group("")
	{
		u.GET("/main_page", au.authPage)
		u.POST("/register", au.register)
		u.POST("/login", au.login)
		u.GET("/logout", au.logout)
	}
}

func (au *AuthHandlers) authPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"block_title": "Test page",
		"URL":         "/v1/auth/login",
	})
}

func (au *AuthHandlers) register(c *gin.Context) {
	var user entity.User
	user.Login = c.PostForm("username")
	user.Password = c.PostForm("password")

	// Validate form input
	if strings.Trim(user.Login, " ") == "" || strings.Trim(user.Password, " ") == "" {
		errorResponse(c, http.StatusBadRequest, "Parameters can't be empty")
		return
	}

	err := au.useCase.AddUser(c, user)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Registration failed")
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered user"})
}

func (au *AuthHandlers) login(c *gin.Context) {

	var user entity.User
	user.Login = c.PostForm("username")
	user.Password = c.PostForm("password")

	// Validate form input
	if strings.Trim(user.Login, " ") == "" || strings.Trim(user.Password, " ") == "" {
		errorResponse(c, http.StatusBadRequest, "Parameters can't be empty")
		return
	}

	ok, err := au.useCase.Login(c, user)
	if err != nil || ok != true {
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
