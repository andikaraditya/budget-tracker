package main

import (
	"github.com/andikaraditya/budget-tracker/backend/internal/source"
	"github.com/andikaraditya/budget-tracker/backend/internal/user"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
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

	app.Listen(":3000")
}
