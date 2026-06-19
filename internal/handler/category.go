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
	protected.Get("/", c.GetByFilter)
	protected.Patch("/:id", c.UpdateById)
	protected.Get("/:id", c.GetById)
	protected.Delete("/:id", c.DeleteById)
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

func (c *CategoryHandler) GetByFilter(ctx *fiber.Ctx) error {
	var req dto.ListCategoriesRequest

	if err := ctx.QueryParser(&req); err != nil {
		log.Println(err)
		return err
	}

	cwt, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	userID := ctx.Locals("user_id").(bson.ObjectID)

	res, err := c.svc.GetByFilter(cwt, userID, req)
	if err != nil {
		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	return ctx.JSON(res)
}

func (c *CategoryHandler) UpdateById(ctx *fiber.Ctx) error {
	var req dto.UpdateCategoryRequest

	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return err
	}

	if err := req.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	cwt, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	id := ctx.Params("id")
	userID := ctx.Locals("user_id").(bson.ObjectID)

	res, err := c.svc.UpdateById(cwt, userID, id, req)
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

func (c *CategoryHandler) GetById(ctx *fiber.Ctx) error {
	cwt, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	id := ctx.Params("id")
	userID := ctx.Locals("user_id").(bson.ObjectID)

	res, err := c.svc.GetById(cwt, userID, id)
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

func (c *CategoryHandler) DeleteById(ctx *fiber.Ctx) error {
	cwt, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	id := ctx.Params("id")
	userID := ctx.Locals("user_id").(bson.ObjectID)

	err := c.svc.DeleteById(cwt, userID, id)
	if err != nil {
		if errors.Is(err, apperror.New("category is linked to one or more transaction")) {
			return ctx.SendStatus(fiber.StatusConflict)
		}

		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.SendStatus(fiber.StatusNotFound)
		}

		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
