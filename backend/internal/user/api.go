package user

import (
	"errors"
	"fmt"

	"github.com/andikaraditya/budget-tracker/backend/internal/api"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(ctx *fiber.Ctx) error {
	req := new(User)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(400).SendString(api.ErrPayload.Error())
	}

	err := Service.CreateUser(*req)
	if err != nil {
		if errors.Is(err, api.ErrPayload) {
			return ctx.Status(400).SendString("email already exists")
		}
		return ctx.Status(500).SendString("internal server error")
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "ok",
	})
}

func Login(ctx *fiber.Ctx) error {
	req := new(User)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	validationError := api.ValidateRequest(req)
	if len(validationError) > 0 {
		return api.SendErrorResponse(ctx, 400, validationError)
	}

	token, err := Service.Login(*req)
	if err != nil {
		fmt.Printf("Error: %v", err)
		if errors.Is(err, api.ErrPayload) {
			return ctx.Status(400).SendString("incorrect password or email")
		}
		return ctx.Status(500).SendString("internal server error")
	}
	return ctx.Status(200).JSON(fiber.Map{
		"token": token,
	})
}