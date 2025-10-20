package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"pura-agung-kertajaya-backend/internal/entity"
	"pura-agung-kertajaya-backend/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	ClearAll()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("rahasia"), bcrypt.DefaultCost)
	user := entity.User{
		Name:     "Admin Test",
		Email:    "admin@puraagungkertajaya.com",
		Password: string(hashed),
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)

	reqBody := model.LoginUserRequest{
		Email:          "admin@puraagungkertajaya.com",
		Password:       "rahasia",
		RecaptchaToken: "dummy_recaptcha_token",
	}
	body, _ := json.Marshal(reqBody)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ := io.ReadAll(response.Body)
	respBody := new(model.WebResponse[model.UserResponse])
	_ = json.Unmarshal(bytes, respBody)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, respBody.Data)
	assert.NotEmpty(t, respBody.Data.Email)
}

func TestCurrentUser(t *testing.T) {
	ClearAll()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("rahasia"), bcrypt.DefaultCost)
	user := entity.User{
		Name:     "Admin Test",
		Email:    "admin@puraagungkertajaya.com",
		Password: string(hashed),
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)

	reqBody := model.LoginUserRequest{
		Email:          "admin@puraagungkertajaya.com",
		Password:       "rahasia",
		RecaptchaToken: "dummy_recaptcha_token",
	}
	body, _ := json.Marshal(reqBody)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("Accept", "application/json")

	loginResp, err := app.Test(loginReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	cookies := loginResp.Cookies()
	assert.NotEmpty(t, cookies)

	var accessToken string
	for _, c := range cookies {
		if c.Name == "access_token" {
			accessToken = c.Value
			break
		}
	}
	assert.NotEmpty(t, accessToken)

	request := httptest.NewRequest(http.MethodGet, "/api/users/_current", nil)
	request.Header.Set("Cookie", "access_token="+accessToken)
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ := io.ReadAll(response.Body)
	resp := new(model.WebResponse[model.UserResponse])
	_ = json.Unmarshal(bytes, resp)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, user.Email, resp.Data.Email)
	assert.NotEmpty(t, resp.Data.CreatedAt)
}

func TestLogout(t *testing.T) {
	ClearAll()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("rahasia"), bcrypt.DefaultCost)
	user := entity.User{
		Name:     "Admin Test",
		Email:    "admin@puraagungkertajaya.com",
		Password: string(hashed),
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)

	reqBody := model.LoginUserRequest{
		Email:          "admin@puraagungkertajaya.com",
		Password:       "rahasia",
		RecaptchaToken: "dummy_recaptcha_token",
	}
	body, _ := json.Marshal(reqBody)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("Accept", "application/json")

	loginResp, err := app.Test(loginReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	cookies := loginResp.Cookies()
	assert.NotEmpty(t, cookies)

	var accessToken string
	for _, c := range cookies {
		if c.Name == "access_token" {
			accessToken = c.Value
			break
		}
	}
	assert.NotEmpty(t, accessToken)

	request := httptest.NewRequest(http.MethodPost, "/api/users/_logout", nil)
	request.Header.Set("Cookie", "access_token="+accessToken)

	response, err := app.Test(request)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestUpdateUser(t *testing.T) {
	ClearAll()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("rahasia"), bcrypt.DefaultCost)
	user := entity.User{
		Name:     "Admin Test",
		Email:    "admin@puraagungkertajaya.com",
		Password: string(hashed),
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)

	reqBody := model.LoginUserRequest{
		Email:          "admin@puraagungkertajaya.com",
		Password:       "rahasia",
		RecaptchaToken: "dummy_recaptcha_token",
	}
	body, _ := json.Marshal(reqBody)

	loginReq := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("Accept", "application/json")

	loginResp, err := app.Test(loginReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	cookies := loginResp.Cookies()
	assert.NotEmpty(t, cookies)
	var accessToken string
	for _, c := range cookies {
		if c.Name == "access_token" {
			accessToken = c.Value
			break
		}
	}
	assert.NotEmpty(t, accessToken)

	updateBody := model.UpdateUserRequest{
		Name: "Administrator Updated",
	}
	updateJSON, _ := json.Marshal(updateBody)

	req := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(updateJSON)))
	req.Header.Set("Cookie", "access_token="+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bytes, _ := io.ReadAll(resp.Body)
	respBody := new(model.WebResponse[model.UserResponse])
	_ = json.Unmarshal(bytes, respBody)

	assert.Equal(t, "Administrator Updated", respBody.Data.Name)

	updated := new(entity.User)
	db.Where("email = ?", "admin@puraagungkertajaya.com").First(updated)
	assert.Equal(t, "Administrator Updated", updated.Name)
}

func TestUpdateUser_InvalidInputs(t *testing.T) {
	ClearAll()

	//  Buat user dummy di DB dengan hash dinamis
	hashed, _ := bcrypt.GenerateFromPassword([]byte("rahasia"), 4) // cost 4 biar cepat
	user := entity.User{
		Name:     "Admin Test",
		Email:    "admin@puraagungkertajaya.com",
		Password: string(hashed),
	}
	_ = db.Create(&user).Error

	//  LOGIN dulu buat dapet cookie
	reqBody := model.LoginUserRequest{
		Email:          "admin@puraagungkertajaya.com",
		Password:       "rahasia",
		RecaptchaToken: "dummy_recaptcha_token",
	}
	body, _ := json.Marshal(reqBody)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(body)))
	loginReq.Header.Set("Content-Type", "application/json")
	loginReq.Header.Set("Accept", "application/json")

	loginResp, err := app.Test(loginReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	var accessToken string
	for _, c := range loginResp.Cookies() {
		if c.Name == "access_token" {
			accessToken = c.Value
		}
	}
	assert.NotEmpty(t, accessToken)

	//  CASE 1: Kosong semua field
	t.Run("Empty Name (should keep old value)", func(t *testing.T) {
		req := model.UpdateUserRequest{Name: ""}
		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(body)))
		request.Header.Set("Cookie", "access_token="+accessToken)
		request.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(request)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// cek apakah name masih sama (tidak berubah)
		bytes, _ := io.ReadAll(resp.Body)
		respBody := new(model.WebResponse[model.UserResponse])
		_ = json.Unmarshal(bytes, respBody)
		assert.Equal(t, "Admin Test", respBody.Data.Name)
	})

	//  CASE 2: Tidak ada cookie (belum login)
	t.Run("Unauthorized (no cookie)", func(t *testing.T) {
		req := model.UpdateUserRequest{Name: "Hacker Attempt"}
		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(body)))
		request.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(request)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	//  CASE 3: Cookie invalid
	t.Run("Invalid Token", func(t *testing.T) {
		req := model.UpdateUserRequest{Name: "Fake User"}
		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(body)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Cookie", "access_token=invalid_token")

		resp, err := app.Test(request)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	//  CASE 4: Update field tidak valid (misal terlalu panjang)
	t.Run("Too Long Name", func(t *testing.T) {
		longName := strings.Repeat("A", 300) // misal validasi max 255
		req := model.UpdateUserRequest{Name: longName}
		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPatch, "/api/users/_current", strings.NewReader(string(body)))
		request.Header.Set("Cookie", "access_token="+accessToken)
		request.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(request)
		assert.Nil(t, err)
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusInternalServerError)
	})
}
