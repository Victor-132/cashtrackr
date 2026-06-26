package main

import (
	"log"

	"github.com/Victor-132/cashtrackr/internal/database"
	"github.com/Victor-132/cashtrackr/internal/handler"
	"github.com/Victor-132/cashtrackr/internal/repository"
	"github.com/Victor-132/cashtrackr/internal/service"
	"github.com/gofiber/fiber/v2"

	reportsHandler "github.com/Victor-132/cashtrackr/internal/reports/handler"
	reportsRepo "github.com/Victor-132/cashtrackr/internal/reports/repository"
	reportsSvc "github.com/Victor-132/cashtrackr/internal/reports/service"
)

func main() {
	app := fiber.New()

	db, err := database.ConnectMongo()
	if err != nil {
		log.Fatal(err)
	}

	database.CreateIndexes(db)

	userRepo := repository.NewUserRepository(db.Collection("users"))
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	userHandler.Register(app)

	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc)
	authHandler.Register(app)

	categoryRepo := repository.NewCategoryRepository(db.Collection("categories"))
	transactionRepo := repository.NewTransactionRepository(db.Collection("transactions"))

	categorySvc := service.NewCategoryService(categoryRepo, transactionRepo)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	categoryHandler.Register(app)

	transactionSvc := service.NewTransactionService(transactionRepo, categoryRepo)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)
	transactionHandler.Register(app)

	reportsTrRepo := reportsRepo.NewTransactionRepository(db.Collection("transactions"))
	reportsTrSvc := reportsSvc.NewTransactionService(reportsTrRepo)
	reportsTrHandler := reportsHandler.NewTransactionHandler(reportsTrSvc)
	reportsTrHandler.Register(app)

	app.Listen(":3000")
}
