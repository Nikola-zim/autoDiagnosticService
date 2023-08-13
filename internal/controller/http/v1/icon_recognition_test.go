package v1

import (
	"autoDiagnosticService/internal/controller/http/v1/mocks"
	"autoDiagnosticService/internal/entity"
	"autoDiagnosticService/pkg/logger"
	"bytes"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const (
	storagePath = "../../../file_storage/images/"
)

func Test_recognition_uploadImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)

	type fields struct {
		useCase     Recognition
		l           logger.Interface
		storagePath string
	}

	type request struct {
		file string
	}

	mockRecognition := mocks.NewMockRecognition(ctrl)
	mockLogger := logger.New("debug")

	tests := []struct {
		name           string
		fields         fields
		request        request
		wantErr        bool
		ExpectStatus   int
		MockUseCaseErr error
	}{
		{
			name: "Test 1",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			request: request{
				file: "./testImg/to_detect234.jpg",
			},
			MockUseCaseErr: nil,
			wantErr:        false,
			ExpectStatus:   http.StatusOK,
		},
		{
			name: "Test 2: resized",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			request: request{
				file: "./testImg/0QAAAgFIq-A-960.jpg",
			},
			MockUseCaseErr: nil,
			wantErr:        false,
			ExpectStatus:   http.StatusOK,
		},
		{
			name: "Test 3: wrong format",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			request: request{
				file: "./testImg/smt_wrong.txt",
			},
			MockUseCaseErr: nil,
			wantErr:        true,
			ExpectStatus:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &recognition{
				useCase:     tt.fields.useCase,
				l:           tt.fields.l,
				storagePath: storagePath,
			}
			// For passing static check
			fmt.Println(r.storagePath)
			// Gin context
			w := httptest.NewRecorder()
			c, handler := gin.CreateTestContext(w)
			// Session for users cookie
			sessionManager := sessions.Sessions("mysession", cookie.NewStore(secret))
			handler.Use(sessionManager)
			// Tested handler
			handler.POST("/recognize", r.uploadImage)

			// File upload
			file, err := os.Open(tt.request.file)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", filepath.Base(tt.request.file))

			if err != nil {
				t.Fatal(err)
			}

			_, err = io.Copy(part, file)
			if err != nil {
				t.Fatal(err)
			}

			err = writer.Close()
			if err != nil {
				t.Fatal(err)
			}

			// Forming request
			c.Request, _ = http.NewRequest("POST", "/recognize", body)
			c.Request.Header.Add("Content-Type", writer.FormDataContentType())

			// Making mock function
			mockRecognition.EXPECT().AddRequest(gomock.Any(), gomock.Any()).Return(tt.MockUseCaseErr).AnyTimes()

			handler.ServeHTTP(w, c.Request)

			// comparing results
			if w.Code != tt.ExpectStatus {
				t.Errorf("expected status %d; got %d", tt.ExpectStatus, w.Code)
			}
		})
	}
}

func Test_recognition_recognizedImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)

	type fields struct {
		useCase     Recognition
		l           logger.Interface
		storagePath string
	}

	mockRecognition := mocks.NewMockRecognition(ctrl)
	mockLogger := logger.New("debug")

	tests := []struct {
		name           string
		fields         fields
		mockResults    []entity.Request
		MockUseCaseErr error
		wantErr        bool
		ExpectStatus   int
	}{
		{
			name: "Test 1: ok",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			MockUseCaseErr: nil,
			mockResults: []entity.Request{
				{
					ResImgPathName: "img1",
				},
				{
					ResImgPathName: "img2",
				},
			},
			wantErr:      false,
			ExpectStatus: http.StatusOK,
		},
		{
			name: "Test 2: no images",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			MockUseCaseErr: nil,
			wantErr:        false,
			ExpectStatus:   http.StatusOK,
		},
		{
			name: "Test 3: internal error",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			MockUseCaseErr: errors.New("some invalid data"),
			wantErr:        true,
			ExpectStatus:   http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &recognition{
				useCase:     tt.fields.useCase,
				l:           tt.fields.l,
				storagePath: storagePath,
			}
			// For passing static check
			fmt.Println(r.storagePath)
			// Gin context
			w := httptest.NewRecorder()
			c, handler := gin.CreateTestContext(w)
			// Session for users cookie
			sessionManager := sessions.Sessions("mysession", cookie.NewStore(secret))
			handler.Use(sessionManager)
			// Tested handler
			handler.GET("/recognizedImages", r.recognizedImages)
			handler.LoadHTMLGlob("../../../controller/static/templates/html/*.html")

			// Forming request
			c.Request, _ = http.NewRequest("GET", "/recognizedImages", nil)

			// Making mock function
			mockRecognition.EXPECT().GetRecognitionAnswersWEB(gomock.Any(), gomock.Any()).Return(tt.mockResults, tt.MockUseCaseErr)

			handler.ServeHTTP(w, c.Request)

			// comparing results
			if w.Code != tt.ExpectStatus {
				t.Errorf("expected status %d; got %d", tt.ExpectStatus, w.Code)
			}
		})
	}
}

func Test_recognition_addPoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)

	type fields struct {
		useCase     Recognition
		l           logger.Interface
		storagePath string
	}
	mockRecognition := mocks.NewMockRecognition(ctrl)
	mockLogger := logger.New("debug")

	tests := []struct {
		name           string
		fields         fields
		message        string
		MockUseCaseErr error
		balanceAdd     int
		wantErr        bool
		ExpectStatus   int
	}{
		{
			name: "Test 1: ok",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			message:        "{ \"points\": 5}",
			balanceAdd:     5,
			MockUseCaseErr: nil,
			wantErr:        false,
			ExpectStatus:   http.StatusOK,
		},
		{
			name: "Test 2: empty message",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			message:        " ",
			MockUseCaseErr: nil,
			wantErr:        false,
			ExpectStatus:   http.StatusBadRequest,
		},
		{
			name: "Test 3: invalid message",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			message:        "{ \"points\": -5}",
			MockUseCaseErr: nil,
			wantErr:        false,
			ExpectStatus:   http.StatusBadRequest,
		},
		{
			name: "Test 4: internal error",
			fields: fields{
				useCase: mockRecognition,
				l:       mockLogger,
			},
			message:        "{ \"points\": -5}",
			balanceAdd:     -5,
			MockUseCaseErr: errors.New("smth wrong with server"),
			wantErr:        true,
			ExpectStatus:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &recognition{
				useCase:     tt.fields.useCase,
				l:           tt.fields.l,
				storagePath: storagePath,
			}
			// For passing static check
			fmt.Println(r.storagePath)
			// Gin context
			w := httptest.NewRecorder()
			c, handler := gin.CreateTestContext(w)
			// Session for users cookie
			sessionManager := sessions.Sessions("mysession", cookie.NewStore(secret))
			handler.Use(sessionManager)
			// Tested handler
			handler.POST("/v1/private/balance/add", r.addPoints)
			handler.LoadHTMLGlob("../../../controller/static/templates/html/*.html")

			// Forming request
			c.Request, _ = http.NewRequest("POST", "/v1/private/balance/add", bytes.NewReader([]byte(tt.message)))

			// Making mock function
			mockRecognition.EXPECT().AddPoints(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.MockUseCaseErr).AnyTimes()

			handler.ServeHTTP(w, c.Request)

			// comparing results
			if w.Code != tt.ExpectStatus {
				t.Errorf("expected status %d; got %d", tt.ExpectStatus, w.Code)
			}
		})
	}
}
