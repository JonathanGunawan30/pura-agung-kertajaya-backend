package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	httpdelivery "pura-agung-kertajaya-backend/internal/delivery/http"
	"pura-agung-kertajaya-backend/internal/model"
	usecasemock "pura-agung-kertajaya-backend/internal/usecase/mock"
)

// setupOrganizationMemberController initializes the controller with mock and fiber app
func setupOrganizationMemberController() (*fiber.App, *usecasemock.OrganizationMemberUsecaseMock) {
	mockUC := &usecasemock.OrganizationMemberUsecaseMock{}
	controller := httpdelivery.NewOrganizationController(mockUC, logrus.New())
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error { // Basic error handler for testing
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return ctx.Status(code).JSON(model.WebResponse[any]{Errors: err.Error()})
		},
	})

	// Register admin routes
	api := app.Group("/api")                      // Assuming a group for API routes
	members := api.Group("/organization-members") // Assuming a group for members
	members.Get("/", controller.GetAll)
	members.Get("/:id", controller.GetByID)
	members.Post("/", controller.Create)
	members.Put("/:id", controller.Update)
	members.Delete("/:id", controller.Delete)

	// Register public route
	public := app.Group("/api/public") // Assuming a group for public routes
	public.Get("/organization-members", controller.GetAllPublic)

	return app, mockUC
}

func TestOrganizationMemberController_GetAllPublic_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()

	// Mock data
	mockResponse := []model.OrganizationResponse{
		{ID: "1", Name: "Member A", Position: "Ketua", PositionOrder: 1, IsActive: true},
		{ID: "2", Name: "Member B", Position: "Sekretaris", PositionOrder: 2, IsActive: true},
	}
	mockUC.On("GetPublic").Return(mockResponse, nil)

	// Create request
	req := httptest.NewRequest("GET", "/api/public/organization-members", nil)

	// Test request
	resp, err := app.Test(req, -1) // Use -1 timeout for simplicity
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Assert response body
	var response model.WebResponse[[]model.OrganizationResponse]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, "Member A", response.Data[0].Name)

	// Assert mock expectations
	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_GetAllPublic_Error(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()

	mockUC.On("GetPublic").Return(([]model.OrganizationResponse)(nil), errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/public/organization-members", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode) // Expecting 500 on generic error

	// Assert mock expectations
	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_GetAll_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()

	mockResponse := []model.OrganizationResponse{
		{ID: "1", Name: "Member A"},
		{ID: "2", Name: "Member B", IsActive: false},
	}
	mockUC.On("GetAll").Return(mockResponse, nil)

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
	app, mockUC := setupOrganizationMemberController()

	mockUC.On("GetAll").Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/organization-members/", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_GetByID_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()
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
	app, mockUC := setupOrganizationMemberController()
	memberID := "notfound"
	// Simulate use case returning a specific "not found" error type if you have one,
	// otherwise a generic error is fine for testing controller's status code mapping.
	mockUC.On("GetByID", memberID).Return((*model.OrganizationResponse)(nil), errors.New("not found")) // Assuming controller maps this to 404

	req := httptest.NewRequest("GET", "/api/organization-members/"+memberID, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode) // Controller should return 404

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Create_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()

	reqBody := model.OrganizationRequest{Name: "New Member", Position: "Anggota", PositionOrder: 5, IsActive: true}
	mockResponse := &model.OrganizationResponse{ID: "newID", Name: "New Member", Position: "Anggota", PositionOrder: 5, IsActive: true}
	mockUC.On("Create", reqBody).Return(mockResponse, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/organization-members/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode) // Expecting 201 Created

	var response model.WebResponse[*model.OrganizationResponse]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "newID", response.Data.ID)
	assert.Equal(t, "New Member", response.Data.Name)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Create_BadBody(t *testing.T) {
	app, _ := setupOrganizationMemberController() // No mock needed

	req := httptest.NewRequest("POST", "/api/organization-members/", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode) // Expecting 400 Bad Request
}

func TestOrganizationMemberController_Create_UsecaseError(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()

	// Simulate validation error from use case
	reqBody := model.OrganizationRequest{Name: "", Position: "Pos"} // Invalid name
	mockUC.On("Create", reqBody).Return((*model.OrganizationResponse)(nil), errors.New("validation failed: Name is required"))

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/organization-members/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode) // Assuming controller maps validation errors to 400

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Update_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()
	memberID := "updateID"
	reqBody := model.OrganizationRequest{Name: "Updated Name", Position: "Ketua", PositionOrder: 1, IsActive: true}
	mockResponse := &model.OrganizationResponse{ID: memberID, Name: "Updated Name", Position: "Ketua", PositionOrder: 1, IsActive: true}

	// Expect Update to be called with ID and Request Body
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
	app, _ := setupOrganizationMemberController()
	memberID := "updateID"

	req := httptest.NewRequest("PUT", "/api/organization-members/"+memberID, bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestOrganizationMemberController_Update_NotFound(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()
	memberID := "notfound"
	reqBody := model.OrganizationRequest{Name: "Update"}

	mockUC.On("Update", memberID, reqBody).Return((*model.OrganizationResponse)(nil), errors.New("not found")) // Simulate use case not found

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/organization-members/"+memberID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	// Decide if your controller returns 404 or 500 on update-not-found
	// Let's assume 500 for now based on your previous controller examples
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode) // Or fiber.StatusNotFound if you map it differently

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Delete_Success(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()
	memberID := "deleteID"

	mockUC.On("Delete", memberID).Return(nil) // Expect Delete to be called with ID

	req := httptest.NewRequest("DELETE", "/api/organization-members/"+memberID, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check response message if your controller sends one
	var response model.WebResponse[any] // Check the type your controller returns
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	// Example assertion: assert.Equal(t, "Organization member deleted successfully", response.Data)

	mockUC.AssertExpectations(t)
}

func TestOrganizationMemberController_Delete_NotFound(t *testing.T) {
	app, mockUC := setupOrganizationMemberController()
	memberID := "notfound"

	mockUC.On("Delete", memberID).Return(errors.New("not found")) // Simulate use case not found

	req := httptest.NewRequest("DELETE", "/api/organization-members/"+memberID, nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	// Decide if your controller returns 404 or 500 on delete-not-found
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode) // Or fiber.StatusNotFound

	mockUC.AssertExpectations(t)
}
