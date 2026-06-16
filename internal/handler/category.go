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

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) CategoryHandler {
	return CategoryHandler{svc}
}

func (c *CategoryHandler) Register(app *fiber.App) {
	protected := app.Group("/v1/categories", middleware.AuthMiddleware)

	protected.Post("/", c.Create)
}

func (c *CategoryHandler) Create(ctx *fiber.Ctx) error {
	var req dto.CreateCategoryRequest

	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return err
	}

	if err := req.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	cwt, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	userID := ctx.Locals("user_id").(bson.ObjectID)

	res, err := c.svc.Create(cwt, userID, req)
	if err != nil {
		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}
