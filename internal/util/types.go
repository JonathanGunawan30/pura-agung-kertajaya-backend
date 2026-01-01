package util

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ImageMap map[string]string

func (m ImageMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *ImageMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte or string failed")
	}

	return json.Unmarshal(bytes, &m)
}

func (ImageMap) GormDataType() string {
	return "json"
}
