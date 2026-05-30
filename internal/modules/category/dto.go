package category

import (
	"github.com/google/uuid"
)

type Response struct {
	ID        uuid.UUID  `json:"id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	IconURL   *string    `json:"icon_url,omitempty"`
	SortOrder int        `json:"sort_order"`
}

func toResponse(c Category) Response {
	return Response{
		ID:        c.ID,
		ParentID:  c.ParentID,
		Name:      c.Name,
		Slug:      c.Slug,
		IconURL:   c.IconURL,
		SortOrder: c.SortOrder,
	}
}

func toResponses(items []Category) []Response {
	out := make([]Response, len(items))
	for i, item := range items {
		out[i] = toResponse(item)
	}
	return out
}
