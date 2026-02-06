package test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (m *UserUsecaseMock) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*model.UserResponse), args.String(1), args.Error(2)
}

func (m *UserUsecaseMock) Current(ctx context.Context, userID int) (*model.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *UserUsecaseMock) UpdateProfile(ctx context.Context, userID int, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *UserUsecaseMock) Logout(ctx context.Context, tokenString string) error {
	args := m.Called(ctx, tokenString)
	return args.Error(0)
}

func setupUserController(t *testing.T) (*fiber.App, *UserUsecaseMock) {
	mockUC := new(UserUsecaseMock)

	cfg := viper.New()
	cfg.Set("cookie.domain", "localhost")

	app, logger, _ := NewTestApp()

	controller := httpdelivery.NewUserController(mockUC, logger, cfg)

	app.Post("/api/users/_login", controller.Login)

	authMiddleware := func(c *fiber.Ctx) error {
		c.Locals("user", &middleware.Auth{ID: 1, Role: "admin"})
		return c.Next()
	}

	app.Get("/api/users/_current", authMiddleware, controller.Current)
	app.Patch("/api/users/_current", authMiddleware, controller.UpdateProfile)
	app.Post("/api/users/_logout", authMiddleware, controller.Logout)

	return app, mockUC
}

func TestUserController_Login_Success(t *testing.T) {
	app, mockUC := setupUserController(t)

	reqBody := model.LoginUserRequest{
		Email:          "admin@puraagungkertajaya.com",
		Password:       "rahasia",
		RecaptchaToken: "dummy",
	}
	expectedUser := &model.UserResponse{ID: 1, Email: reqBody.Email, Name: "Admin"}
	expectedToken := "mock-jwt-token"

	mockUC.On("Login", mock.Anything, mock.AnythingOfType("*model.LoginUserRequest")).
		Return(expectedUser, expectedToken, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookies := resp.Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "access_token" && c.Value == expectedToken {
			found = true
			break
		}
	}
	assert.True(t, found, "Access token cookie should be set")
}

func TestUserController_Login_Fail(t *testing.T) {
	app, mockUC := setupUserController(t)

	reqBody := model.LoginUserRequest{
		Email:    "wrong@example.com",
		Password: "wrong",
	}

	mockUC.On("Login", mock.Anything, mock.AnythingOfType("*model.LoginUserRequest")).
		Return((*model.UserResponse)(nil), "", errors.New("invalid credentials"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUserController_Current_Success(t *testing.T) {
	app, mockUC := setupUserController(t)

	expectedUser := &model.UserResponse{
		ID:    1,
		Name:  "Admin Test",
		Email: "admin@puraagungkertajaya.com",
	}

	mockUC.On("Current", mock.Anything, 1).Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bytes, _ := io.ReadAll(resp.Body)
	var webResp model.WebResponse[model.UserResponse]
	json.Unmarshal(bytes, &webResp)

	assert.Equal(t, expectedUser.Email, webResp.Data.Email)
}

func TestUserController_UpdateProfile_Success(t *testing.T) {
	app, mockUC := setupUserController(t)

	reqBody := model.UpdateUserRequest{Name: "New Name"}
	expectedUser := &model.UserResponse{ID: 1, Name: "New Name"}

	mockUC.On("UpdateProfile", mock.Anything, 1, mock.AnythingOfType("*model.UpdateUserRequest")).
		Return(expectedUser, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserController_Logout_Success(t *testing.T) {
	app, mockUC := setupUserController(t)

	tokenString := "valid-token"

	mockUC.On("Logout", mock.Anything, tokenString).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/users/_logout", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: tokenString})

	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	cookies := resp.Cookies()
	cleared := false
	for _, c := range cookies {
		if c.Name == "access_token" && (c.MaxAge < 0 || c.Value == "") {
			cleared = true
			break
		}
	}
	assert.True(t, cleared, "Access token cookie should be cleared")
}
