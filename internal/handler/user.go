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

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return UserHandler{svc}
}

func (u *UserHandler) Register(app *fiber.App) {
	pub := app.Group("/v1/users")

	pub.Post("", u.Create)

	protected := pub.Group("/", middleware.AuthMiddleware)

	protected.Patch("/me", u.UpdateProfile)
	protected.Patch("/me/password", u.ChangePassword)
}

func (u *UserHandler) Create(ctx *fiber.Ctx) error {
	var req dto.CreateUserRequest

	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return err
	}

	if err := req.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	res, err := u.svc.Create(c, req)
	if err != nil {
		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (u *UserHandler) UpdateProfile(ctx *fiber.Ctx) error {
	var req dto.UpdateUserRequest

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

	res, err := u.svc.UpdateProfile(c, userID, req)
	if err != nil {
		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	return ctx.JSON(res)
}

func (u *UserHandler) ChangePassword(ctx *fiber.Ctx) error {
	var req dto.ChangePasswordRequest

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

	if err := u.svc.ChangePassword(c, userID, req); err != nil {
		var appError apperror.AppError
		if errors.As(err, &appError) {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	return nil
}
