package userhandler

import (
	"bytes"
	"film_library/internal/domains"
	"film_library/internal/repositories/postgres/userrepo"
	mock_services "film_library/internal/services/mocks"
	userservice "film_library/internal/services/userservice"
	"film_library/pkg/mux"
	"film_library/pkg/validation"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserHandlerRegister(t *testing.T) {
	type mockBehavior func(r *mock_services.MockUserService, user domains.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domains.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Correct",
			inputBody: `{"login":"denis", "password":"password","role":"viewer"}`,
			inputUser: domains.User{Login: "denis", Password: "password", Role: "viewer"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().CreateUser(user).Return("token", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:      "Invalid role",
			inputBody: `{"login":"denis", "password":"password","role":"aboba"}`,
			inputUser: domains.User{Login: "denis", Password: "password", Role: "aboba"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().CreateUser(user).Return("", &validation.ValidateError{fmt.Errorf("invalid role")})
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"errors":["invalid role"]}`,
		},
		{
			name:      "Invalid credentials",
			inputBody: `{"login":"", "password":"1","role":"aboba"}`,
			inputUser: domains.User{Login: "", Password: "1", Role: "aboba"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().CreateUser(user).Return("",
					&validation.ValidateError{
						fmt.Errorf("invalid login length"),
						fmt.Errorf("invalid password length"),
						fmt.Errorf("invalid role"),
					},
				)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"errors":["invalid login length","invalid password length","invalid role"]}`,
		},
		{
			name:                 "Json unmarshal error",
			inputBody:            ``,
			mockBehavior:         func(r *mock_services.MockUserService, user domains.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad request"}`,
		},
		{
			name:                 "Json unmarshal error",
			inputBody:            ``,
			mockBehavior:         func(r *mock_services.MockUserService, user domains.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad request"}`,
		},
		{
			name:      "Invalid credentials",
			inputBody: `{"login":"admin", "password":"password","role":"admin"}`,
			inputUser: domains.User{Login: "admin", Password: "password", Role: "admin"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().CreateUser(user).Return("", userrepo.ErrAlreadyExists)
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: `{"error":"user already exists"}`,
		},
		{
			name:      "Unknown error",
			inputBody: `{"login":"123", "password":"123","role":"123"}`,
			inputUser: domains.User{Login: "123", Password: "123", Role: "123"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().CreateUser(user).Return("", fmt.Errorf("some error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"unknown error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_services.NewMockUserService(c)
			handler := UserHandler{service: service}
			tc.mockBehavior(service, tc.inputUser)

			r := mux.New()
			r.HandleFunc("POST /register", handler.Register)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			if tc.expectedStatusCode != w.Code {
				t.Errorf("expected: %d\ngot: %d", tc.expectedStatusCode, w.Code)
			}

			if tc.expectedResponseBody != w.Body.String() {
				t.Errorf("expected: %s\ngot: %s", tc.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestUserHandlerLogin(t *testing.T) {
	type mockBehavior func(r *mock_services.MockUserService, user domains.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domains.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Correct",
			inputBody: `{"login":"denis", "password":"password","role":"viewer"}`,
			inputUser: domains.User{Login: "denis", Password: "password", Role: "viewer"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().Login(user.Login, user.Password).Return("token", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:      "Invalid password",
			inputBody: `{"login":"denis", "password":"1","role":"viewer"}`,
			inputUser: domains.User{Login: "denis", Password: "1", Role: "viewer"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().Login(user.Login, user.Password).Return("", userservice.ErrInvalidPassword)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"invalid login or password"}`,
		},
		{
			name:      "User not found",
			inputBody: `{"login":"adfgfdag", "password":"1","role":"viewer"}`,
			inputUser: domains.User{Login: "adfgfdag", Password: "1", Role: "viewer"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().Login(user.Login, user.Password).Return("", userservice.ErrNotFound)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"invalid login or password"}`,
		},
		{
			name:                 "Json unmarshal error",
			inputBody:            ``,
			mockBehavior:         func(r *mock_services.MockUserService, user domains.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"bad request"}`,
		},
		{
			name:      "Unknown error",
			inputBody: `{"login":"123", "password":"123","role":"123"}`,
			inputUser: domains.User{Login: "123", Password: "123", Role: "123"},
			mockBehavior: func(r *mock_services.MockUserService, user domains.User) {
				r.EXPECT().Login(user.Login, user.Password).Return("", fmt.Errorf("some error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"unknown error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_services.NewMockUserService(c)
			handler := UserHandler{service: service}
			tc.mockBehavior(service, tc.inputUser)

			r := mux.New()
			r.HandleFunc("POST /login", handler.Login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			if tc.expectedStatusCode != w.Code {
				t.Errorf("expected: %d\ngot: %d", tc.expectedStatusCode, w.Code)
			}

			if tc.expectedResponseBody != w.Body.String() {
				t.Errorf("expected: %s\ngot: %s", tc.expectedResponseBody, w.Body.String())
			}
		})
	}
}
