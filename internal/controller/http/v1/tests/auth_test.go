package tests_test

import (
	v1 "autoDiagnosticService/internal/controller/http/v1"
	"autoDiagnosticService/internal/controller/http/v1/mocks"
	"autoDiagnosticService/internal/entity"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	defaultUsername = "username"
	defaultPassword = "password"
)

var secret = []byte("secret")

func TestLoginWEB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		Name           string
		Username       string
		Password       string
		User           entity.User
		MockUseCaseOK  bool
		MockUseCaseErr error
		ExpectStatus   int
	}{
		{
			Name:           "Successful LoginWEB",
			Username:       defaultUsername,
			Password:       defaultPassword,
			User:           entity.User{Login: defaultUsername, Password: defaultPassword},
			MockUseCaseOK:  true,
			MockUseCaseErr: nil,
			ExpectStatus:   http.StatusOK,
		},
		{
			Name:           "Empty Parameters",
			Username:       "  ",
			Password:       "  ",
			User:           entity.User{Login: "  ", Password: "  "},
			MockUseCaseOK:  false,
			MockUseCaseErr: nil,
			ExpectStatus:   http.StatusBadRequest,
		},
		{
			Name:           "Invalid Credentials",
			Username:       "qwerqwerqwer",
			Password:       "qwerqwerqwer",
			User:           entity.User{Login: "qwerqwerqwer", Password: "qwerqwerqwer"},
			MockUseCaseOK:  false,
			MockUseCaseErr: nil,
			ExpectStatus:   http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			loginURL := "/v1/auth/login"
			w := httptest.NewRecorder()

			c, r := gin.CreateTestContext(w)
			mockUseCase := mocks.NewMockRecognition(ctrl)
			au := &v1.AuthHandlers{
				UseCase: mockUseCase,
			}
			sessionManager := sessions.Sessions("mysession", cookie.NewStore(secret))
			r.Use(sessionManager)
			r.POST(loginURL, au.LoginWEB)

			data := url.Values{}
			data.Add("username", tc.Username)
			data.Add("password", tc.Password)
			body := strings.NewReader(data.Encode())
			c.Request, _ = http.NewRequest(http.MethodPost, loginURL, body)
			c.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			mockUseCase.EXPECT().Login(gomock.Any(), tc.User).Return(tc.MockUseCaseOK, tc.MockUseCaseErr).AnyTimes()
			r.ServeHTTP(w, c.Request)

			if w.Code != tc.ExpectStatus {
				t.Errorf("expected status %d; got %d", tc.ExpectStatus, w.Code)
			}
		})
	}
}

func TestAuthHandlers_register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		Name           string
		Username       string
		Password       string
		User           entity.User
		MockUseCaseOK  bool
		MockUseCaseErr error
		ExpectStatus   int
	}{
		{
			Name:           "Successful register",
			Username:       defaultUsername,
			Password:       defaultPassword,
			User:           entity.User{Login: defaultUsername, Password: defaultPassword},
			MockUseCaseErr: nil,
			ExpectStatus:   http.StatusOK,
		},
		{
			Name:           "Empty Parameters",
			Username:       "  ",
			Password:       "  ",
			User:           entity.User{Login: "  ", Password: "  "},
			MockUseCaseErr: nil,
			ExpectStatus:   http.StatusBadRequest,
		},
		{
			Name:           "Invalid Credentials",
			Username:       "#@!$#",
			Password:       "qwerqwerqwer",
			User:           entity.User{Login: "#@!$#", Password: "qwerqwerqwer"},
			MockUseCaseErr: nil,
			ExpectStatus:   http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			loginURL := "/v1/auth/register"
			w := httptest.NewRecorder()

			c, r := gin.CreateTestContext(w)
			mockUseCase := mocks.NewMockRecognition(ctrl)
			au := &v1.AuthHandlers{
				UseCase: mockUseCase,
			}
			sessionManager := sessions.Sessions("mysession", cookie.NewStore(secret))
			r.Use(sessionManager)
			r.POST(loginURL, au.Register)

			data := url.Values{}
			data.Add("username", tc.Username)
			data.Add("password", tc.Password)
			body := strings.NewReader(data.Encode())
			c.Request, _ = http.NewRequest(http.MethodPost, loginURL, body)
			c.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			mockUseCase.EXPECT().AddUser(gomock.Any(), tc.User).Return(tc.MockUseCaseErr).AnyTimes()
			r.ServeHTTP(w, c.Request)

			if w.Code != tc.ExpectStatus {
				t.Errorf("expected status %d; got %d", tc.ExpectStatus, w.Code)
			}
		})
	}
}
