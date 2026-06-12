package auth

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetUserID(ctx *fiber.Ctx) bson.ObjectID {
	return ctx.Locals("user_id").(bson.ObjectID)

}
