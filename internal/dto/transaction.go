package dto

import (
	"strings"
	"time"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
)

type CreateTransactionRequest struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Amount          int       `json:"amount"`
	Type            string    `json:"type"`
	TransactionDate time.Time `json:"transaction_date"`
}

func (c *CreateTransactionRequest) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return apperror.New("invalid title")
	}

	if c.Amount <= 0 {
		return apperror.New("invalid amount")
	}

	switch c.Type {
	case "income", "expense":
	default:
		return apperror.New("invalid type")
	}

	if c.TransactionDate.IsZero() {
		return apperror.New("invalid date")
	}

	return nil
}

type TransactionResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	Amount          int       `json:"amount"`
	Type            string    `json:"type"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
}

type ListTransactionsRequest struct {
	Page      int        `query:"page"`
	Limit     int        `query:"limit"`
	Type      string     `query:"type"`
	StartDate *time.Time `query:"start_date"`
	EndDate   *time.Time `query:"end_date"`
}

type ListTransactionsResponse struct {
	Data       []TransactionResponse `json:"data"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalItems int                   `json:"total_items"`
	TotalPages int                   `json:"total_pages"`
}
