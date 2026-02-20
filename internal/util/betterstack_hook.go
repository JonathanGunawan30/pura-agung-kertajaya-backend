package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type BetterStackHook struct {
	Token string
	URL   string
}

func (h *BetterStackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *BetterStackHook) Fire(entry *logrus.Entry) error {
	payload := map[string]any{
		"dt":      entry.Time.Format(time.RFC3339),
		"level":   entry.Level.String(),
		"message": entry.Message,
		"service": "pura-agung-kertajaya-backend",
	}

	for k, v := range entry.Data {
		payload[k] = v
	}

	jsonPayload, _ := json.Marshal(payload)

	go func(data []byte) {
		req, err := http.NewRequest(http.MethodPost, h.URL, bytes.NewBuffer(data))
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+h.Token)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
		}
	}(jsonPayload)
	return nil
}
