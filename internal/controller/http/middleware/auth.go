package middleware

import (
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

const userKey = "user"

type Auth struct {
	useCase usecase.Recognition
	l       logger.Interface
	path    string
	nextURL string
}

func NewAuth(l logger.Interface, useCase usecase.Recognition) *Auth {
	return &Auth{
		l:       l,
		useCase: useCase,
	}
}

func (au *Auth) NewAuthRoutes(handler *gin.RouterGroup, path string, nextURL string) {
	au.path = path
	au.nextURL = nextURL
	// Login and logout routes
	u := handler.Group("")
	{
		u.GET("/main_page", au.authPage)
		u.POST("/register", au.register)
		u.POST("/login", au.login)
		u.GET("/logout", au.logout)
	}
}

func (au *Auth) authPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"block_title": "Test page",
		"URL":         au.path + "/login",
	})
}

func (au *Auth) register(c *gin.Context) {
	var user entity.User
	user.Login = c.PostForm("uname")
	user.Password = c.PostForm("psw")

	au.l.Info("Username: %s; Password: %s", user.Login, user.Password)

	err := au.useCase.AddUser(c, user)
	if err != nil {
		c.HTML(http.StatusOK, "auth.html", gin.H{
			"block_title": "Authorization",
			"status":      "Registration failed!",
		})
		return
	}
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"block_title": "Authorization",
		"status":      "User added! Please, login",
	})
}

func (au *Auth) login(c *gin.Context) {

	var user entity.User
	user.Login = c.PostForm("uname")
	user.Password = c.PostForm("psw")

	au.l.Info("Username: %s; Password: %s", user.Login, user.Password)

	ok, err := au.useCase.Login(c, user)
	if err != nil || ok != true {
		c.HTML(http.StatusUnauthorized, "auth.html", gin.H{
			"block_title": "Authorization",
			"status":      "authentication failed!",
		})
	} else {
		//Устанавливаем куку
		c.SetCookie("user_id", "cookie", 3600, "/", "localhost", false, true)
		// переброска данных далее в запрос
		c.Set("user_ID", "cookie")
		// Передача запроса в handler
		au.l.Info("neeeext")
		c.Redirect(http.StatusFound, au.nextURL)
	}
}

func (au *Auth) logout(c *gin.Context) {

	var user entity.User
	user.Login = c.PostForm("uname")
	user.Password = c.PostForm("psw")

	au.l.Info("Username: %s; Password: %s", user.Login, user.Password)

	ok, err := au.useCase.Login(c, user)
	if err != nil || ok != true {
		c.HTML(http.StatusOK, "auth.html", gin.H{
			"block_title": "Authorization",
			"status":      "authentication failed!",
		})
	} else {
		// Передача запроса в handler
		c.Next()
	}
}

// AuthRequired - middleware для cookie
func (au *Auth) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, err := c.Cookie("user_id")
		if err != nil {
			// Redirect to auth page
			c.Redirect(http.StatusTemporaryRedirect, au.path+"/main_page")
			return
		}

		//log.Printf("user_ID: %s \n", cookie)
		au.l.Info("user_ID: %v \n", user_id)
		// Передача запроса в handler
		c.Next()
	}
}
