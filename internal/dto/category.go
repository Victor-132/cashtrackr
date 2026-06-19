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

type UpdateCategoryRequest struct {
	Name *string `json:"name"`
}

func (u *UpdateCategoryRequest) Validate() error {
	if u.Name != nil && strings.TrimSpace(*u.Name) == "" {
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

type ListCategoriesRequest struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type ListCategoriesResponse struct {
	Data       []CategoryResponse `json:"data"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalItems int                `json:"total_items"`
	TotalPages int                `json:"total_pages"`
}
