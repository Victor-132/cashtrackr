package handler

import (
	"context"
	"errors"
	"log"
	"time"

	apperror "github.com/Victor-132/cashtrackr/internal/app_error"
	"github.com/Victor-132/cashtrackr/internal/dto"
	"github.com/Victor-132/cashtrackr/internal/middleware"
	"github.com/Victor-132/cashtrackr/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionHandler struct {
	svc service.TrasnactionService
}

func NewTransactionHandler(svc service.TrasnactionService) TransactionHandler {
	return TransactionHandler{svc}
}

func (t *TransactionHandler) Register(app *fiber.App) {
	protected := app.Group("/v1/transactions", middleware.AuthMiddleware)

	protected.Post("/", t.Create)
	protected.Get("/", t.GetByFilter)
	protected.Get("/:id", t.GetById)
}

func (t *TransactionHandler) Create(ctx *fiber.Ctx) error {
	var req dto.CreateTransactionRequest

	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return err
	}

	if err := req.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	userID := ctx.Locals("user_id").(bson.ObjectID)

	res, err := t.svc.Create(c, userID, req)
	if err != nil {
		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (t *TransactionHandler) GetByFilter(ctx *fiber.Ctx) error {
	var req dto.ListTransactionsRequest

	if err := ctx.QueryParser(&req); err != nil {
		log.Println(err)
		return err
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	userID := ctx.Locals("user_id").(bson.ObjectID)

	res, err := t.svc.GetByFilter(c, userID, req)
	if err != nil {
		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	return ctx.JSON(res)
}

func (t *TransactionHandler) GetById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	userID := ctx.Locals("user_id").(bson.ObjectID)

	res, err := t.svc.GetById(c, userID, id)
	if err != nil {
		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	if res == nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(res)
}
