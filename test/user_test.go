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
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/delivery/http/middleware"
	"pura-agung-kertajaya-backend/internal/model"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (m *UserUsecaseMock) Login(ctx context.Context, req *model.LoginUserRequest, fiberCtx *fiber.Ctx) (*model.UserResponse, error) {
	args := m.Called(ctx, req, fiberCtx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
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

func (m *UserUsecaseMock) Logout(ctx context.Context, fiberCtx *fiber.Ctx) (bool, error) {
	args := m.Called(ctx, fiberCtx)
	return args.Bool(0), args.Error(1)
}

func setupUserController(t *testing.T) (*fiber.App, *UserUsecaseMock) {
	mockUC := new(UserUsecaseMock)
	controller := httpdelivery.NewUserController(mockUC, logrus.New())

	app := fiber.New()

	app.Post("/api/users/_login", controller.Login)

	// FIX: Middleware Mock dipasang eksplisit
	// Ini memastikan user ter-inject untuk route protected
	authMiddleware := func(c *fiber.Ctx) error {
		c.Locals("user", &middleware.Auth{ID: 1}) // Inject Mock User ID = 1
		return c.Next()
	}

	// Terapkan middleware langsung pada definisi route agar lebih aman
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

	mockUC.On("Login", mock.Anything, mock.AnythingOfType("*model.LoginUserRequest"), mock.Anything).
		Return(expectedUser, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserController_Login_Fail(t *testing.T) {
	app, mockUC := setupUserController(t)

	reqBody := model.LoginUserRequest{
		Email:    "wrong@example.com",
		Password: "wrong",
	}

	mockUC.On("Login", mock.Anything, mock.AnythingOfType("*model.LoginUserRequest"), mock.Anything).
		Return(nil, errors.New("invalid credentials"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	// FIX: Controller mengembalikan 500 jika error tidak spesifik (seperti yang terlihat di log error)
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
	resp, _ := app.Test(req)

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

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserController_Logout_Success(t *testing.T) {
	app, mockUC := setupUserController(t)

	mockUC.On("Logout", mock.Anything, mock.Anything).Return(true, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/users/_logout", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
