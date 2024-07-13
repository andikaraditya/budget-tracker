package main

import (
	"github.com/andikaraditya/budget-tracker/backend/internal/category"
	"github.com/andikaraditya/budget-tracker/backend/internal/record"
	"github.com/andikaraditya/budget-tracker/backend/internal/source"
	"github.com/andikaraditya/budget-tracker/backend/internal/transfer"
	"github.com/andikaraditya/budget-tracker/backend/internal/user"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/uuid", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"uuid": uuid.NewString(),
		})
	})

	app.Post("/register", user.CreateUser)
	app.Post("login", user.Login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	app.Post("/sources", source.CreateSource)
	app.Get("/sources", source.GetSources)
	app.Get("/sources/:sourceId", source.GetSource)
	app.Put("/sources/:sourceId", source.UpdateSource)

	app.Post("/categories", category.CreateCategory)
	app.Get("/categories", category.GetCategories)
	app.Get("/categories/:categoryId", category.GetCategory)
	app.Put("/categories/:categoryId", category.UpdateCategory)

	app.Post("/transfers", transfer.CreateTransfer)
	app.Get("/transfers", transfer.GetTransfers)
	app.Get("/transfers/:transferId", transfer.GetTransfer)
	app.Put("/transfers/:transferId", transfer.UpdateTransfer)

	app.Post("/records", record.CreateRecord)
	app.Get("/records", record.GetRecords)
	app.Get("/records/:recordId", record.GetRecord)
	app.Put("/records/:recordId", record.UpdateRecord)

	app.Listen(":3000")
}
