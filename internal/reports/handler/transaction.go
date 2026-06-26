package handler

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Victor-132/cashtrackr/internal/middleware"
	"github.com/Victor-132/cashtrackr/internal/reports/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Transaction struct {
	svc service.Transaction
}

func NewTransactionHandler(svc service.Transaction) Transaction {
	return Transaction{svc}
}

func (t *Transaction) Register(app *fiber.App) {
	protected := app.Group("/v1/reports", middleware.AuthMiddleware)

	protected.Get("/monthly-summary", t.GetMonthlySummary)
	protected.Get("/expenses-by-category", t.GetExpensesByCategory)
	protected.Get("/monthly-evolution", t.GetMonthlyEvolution)
}

func (t *Transaction) GetMonthlySummary(ctx *fiber.Ctx) error {
	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid year")
	}

	month, err := strconv.Atoi(ctx.Query("month"))
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid month")
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	userID := ctx.Locals("user_id").(bson.ObjectID)

	result, err := t.svc.GetMonthlySummary(c, userID, year, month)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (t *Transaction) GetExpensesByCategory(ctx *fiber.Ctx) error {
	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid year")
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	userID := ctx.Locals("user_id").(bson.ObjectID)

	result, err := t.svc.GetExpensesByCategory(c, userID, year)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (t *Transaction) GetMonthlyEvolution(ctx *fiber.Ctx) error {
	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid year")
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	userID := ctx.Locals("user_id").(bson.ObjectID)

	result, err := t.svc.GetMonthlyEvolution(c, userID, year)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
