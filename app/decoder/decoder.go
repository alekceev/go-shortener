package decoder

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

func Decode(i uuid.UUID) (string, error) {
	return base64.RawURLEncoding.EncodeToString([]byte(i.String())), nil
}

func Encode(s string) (uuid.UUID, error) {
	id, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unexpected character in str: %v", err)
	}

	u, _ := uuid.Parse(string(id))
	return u, nil
}
