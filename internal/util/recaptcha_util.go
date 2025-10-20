package util

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

type RecaptchaUtil struct {
	SecretKey string
	Env       string
	Client    *http.Client
}

func NewRecaptchaUtil(v *viper.Viper) *RecaptchaUtil {
	return &RecaptchaUtil{
		SecretKey: v.GetString("recaptcha.secret"),
		Env:       v.GetString("app.env"),
		Client:    &http.Client{},
	}
}

func (r *RecaptchaUtil) Verify(ctx context.Context, token string) bool {
	if r.Env == "development" {
		return true
	}

	data := map[string]string{
		"secret":   r.SecretKey,
		"response": token,
	}

	body, _ := json.Marshal(data)
	resp, err := r.Client.Post("https://www.google.com/recaptcha/api/siteverify",
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	return result.Success
}
