package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

func setupOrganizationMemberController(t *testing.T) (*fiber.App, *usecasemock.OrganizationMemberUsecaseMock) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	controller := httpdelivery.NewOrganizationController(mockUC, logrus.New())
	app := fiber.New(fiber.Config{
		StrictRouting: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "An internal server error occurred. Please try again later."
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				code = fiber.StatusNotFound
				message = "The requested resource was not found."
			} else if _, ok := err.(validator.ValidationErrors); ok {
				code = fiber.StatusBadRequest
				message = "Validation failed"
			}
			return ctx.Status(code).JSON(model.WebResponse[any]{Errors: message})
		},
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("entity_type", "pura")
		return c.Next()
	})

	api := app.Group("/api")
	members := api.Group("/organization-members")

	members.Get("/", controller.GetAll)
	members.Get("/:id", controller.GetByID)
	members.Post("/", controller.Create)
	members.Put("/:id", controller.Update)
	members.Delete("/:id", controller.Delete)

	public := app.Group("/api/public")
	public.Get("/organization-members", controller.GetAllPublic)

	return app, mockUC
}

func TestOrganizationMemberController_GetAllPublic_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)

	items := []model.OrganizationResponse{
		{ID: "1", Name: "Member A", Position: "Ketua", PositionOrder: 1, IsActive: true},
		{ID: "2", Name: "Member B", Position: "Sekretaris", PositionOrder: 2, IsActive: true},
	}
	mockUC.On("GetPublic", "").Return(items, nil)

	req := httptest.NewRequest("GET", "/api/public/organization-members", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.OrganizationResponse]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, "Member A", response.Data[0].Name)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_GetAllPublic_Error(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)

	mockUC.On("GetPublic", "").Return(([]model.OrganizationResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/organization-members", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_GetAll_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)

	mockResponse := []model.OrganizationResponse{
		{ID: "1", Name: "Member A"},
		{ID: "2", Name: "Member B", IsActive: false},
	}
	mockUC.On("GetAll", "pura").Return(mockResponse, nil)

	req := httptest.NewRequest("GET", "/api/organization-members/", nil) // Note trailing slash might matter depending on StrictRouting

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[[]model.OrganizationResponse]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Data, 2)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_GetAll_Error(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)

	mockUC.On("GetAll", "pura").Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/organization-members/", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_GetByID_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)
	memberID := "member123"
	mockResponse := &model.OrganizationResponse{ID: memberID, Name: "Member Found"}
	mockUC.On("GetByID", memberID).Return(mockResponse, nil)

	req := httptest.NewRequest("GET", "/api/organization-members/"+memberID, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationResponse]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, memberID, response.Data.ID)
	assert.Equal(t, "Member Found", response.Data.Name)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_GetByID_NotFound(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)
	memberID := "notfound"
	mockUC.On("GetByID", memberID).Return((*model.OrganizationResponse)(nil), gorm.ErrRecordNotFound) // Assuming controller maps this to 404

	req := httptest.NewRequest("GET", "/api/organization-members/"+memberID, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Create_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)

	reqBody := model.CreateOrganizationRequest{
		Name:          "New Member",
		Position:      "Anggota",
		PositionOrder: 5,
		IsActive:      true,
	}
	mockResponse := &model.OrganizationResponse{
		ID:            "newID",
		Name:          "New Member",
		Position:      "Anggota",
		PositionOrder: 5,
		IsActive:      true,
	}
	mockUC.On("Create", "pura", reqBody).Return(mockResponse, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/organization-members/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationResponse]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "newID", response.Data.ID)
	assert.Equal(t, "New Member", response.Data.Name)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Create_BadBody(t *testing.T) {
	app, _ := setupOrganizationMemberController(t)

	req := httptest.NewRequest("POST", "/api/organization-members/", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestOrganizationMemberController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)

	reqBody := model.CreateOrganizationRequest{Name: "", Position: "Pos"}
	// entity_type is set to "pura" by test middleware
	validate := validator.New()
	err := validate.Struct(reqBody)
	var validationErrs validator.ValidationErrors
	errors.As(err, &validationErrs)
	mockUC.On("Create", "pura", reqBody).Return((*model.OrganizationResponse)(nil), validationErrs)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/organization-members/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Update_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)
	memberID := "updateID"
	reqBody := model.UpdateOrganizationRequest{
		Name:          "Updated Name",
		Position:      "Ketua",
		PositionOrder: 1,
		IsActive:      true,
	}
	mockResponse := &model.OrganizationResponse{
		ID:            memberID,
		Name:          "Updated Name",
		Position:      "Ketua",
		PositionOrder: 1,
		IsActive:      true,
	}

	mockUC.On("Update", memberID, reqBody).Return(mockResponse, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-members/"+memberID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[*model.OrganizationResponse]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, memberID, response.Data.ID)
	assert.Equal(t, "Updated Name", response.Data.Name)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Update_BadBody(t *testing.T) {
	app, _ := setupOrganizationMemberController(t)
	memberID := "updateID"

	req := httptest.NewRequest("PUT", "/api/organization-members/"+memberID, bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestOrganizationMemberController_Update_NotFound(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)
	memberID := "notfound"
	reqBody := model.UpdateOrganizationRequest{Name: "Update"}

	mockUC.On("Update", memberID, reqBody).Return((*model.OrganizationResponse)(nil), gorm.ErrRecordNotFound) // Simulate use case not found

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-members/"+memberID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Delete_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)
	memberID := "deleteID"

	mockUC.On("Delete", memberID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/organization-members/"+memberID, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response model.WebResponse[any]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Delete_NotFound(t *testing.T) {
	app, mockUC := setupOrganizationMemberController(t)
	memberID := "notfound"

	mockUC.On("Delete", memberID).Return(gorm.ErrRecordNotFound)

	req := httptest.NewRequest("DELETE", "/api/organization-members/"+memberID, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUC.AssertExpectations(t)
}
