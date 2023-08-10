package v1

import (
	"autoDiagnosticService/internal/controller/http/v1/mocks"
	"autoDiagnosticService/pkg/logger"
	"bytes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
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
