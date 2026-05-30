package campus

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Domain    *string   `json:"domain,omitempty"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

func toResponse(c Campus) Response {
	return Response{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		Domain:    c.Domain,
		Country:   c.Country,
		City:      c.City,
		CreatedAt: c.CreatedAt,
	}
}

func toResponses(items []Campus) []Response {
	out := make([]Response, len(items))
	for i, item := range items {
		out[i] = toResponse(item)
	}
	return out
}
