package campus

import (
	"fmt"

	"github.com/google/uuid"
)

func parseUUID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid id")
	}
	return id, nil
}
