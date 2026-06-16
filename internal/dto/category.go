package dto

import (
	"strings"
	"time"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (c *CreateCategoryRequest) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return apperror.New("invalid name")
	}

	return nil
}

type CategoryResponse struct {
	ID        string    `json:"_id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
