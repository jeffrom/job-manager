package resource

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type NullRawMessage struct {
	json.RawMessage
	Valid bool
}

// Scan implements the Scanner interface.
func (n *NullRawMessage) Scan(value interface{}) error {
	if value == nil {
		n.RawMessage, n.Valid = json.RawMessage{}, false
		return nil
	}
	buf, ok := value.([]byte)

	if !ok {
		return errors.New("canot parse to bytes")
	}

	n.RawMessage, n.Valid = buf, true

	return nil
}

// Value implements the driver Valuer interface.
func (n NullRawMessage) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.RawMessage, nil
}
