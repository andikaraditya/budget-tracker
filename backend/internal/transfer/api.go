package transfer

import (
	"errors"
	"fmt"

	"github.com/andikaraditya/budget-tracker/backend/internal/api"
	apiParams "github.com/andikaraditya/budget-tracker/backend/internal/params"
	"github.com/gofiber/fiber/v2"
)

func CreateTransfer(c *fiber.Ctx) error {
	req := new(Transfer)

	req.UserId = api.GetUserId(c)

	if err := c.BodyParser(req); err != nil {
		return c.Status(500).SendString("internal server error")
	}

	validationError := api.ValidateRequest(req)
	if len(validationError) > 0 {
		return api.SendErrorResponse(c, 400, validationError)
	}

	req.UserId = api.GetUserId(c)

	err := Service.createTransfer(req)
	if err != nil {
		return c.Status(500).SendString("internal server error")
	}
	return c.Status(201).JSON(req)
}

func GetTransfer(c *fiber.Ctx) error {
	req := new(Transfer)

	req.ID = c.Params("transferId")

	userId := api.GetUserId(c)

	if err := Service.getTransfer(req, userId); err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return api.SendErrorResponse(c, 404, "transfer not found")
		}
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	return c.Status(200).JSON(req)
}

func GetTransfers(c *fiber.Ctx) error {
	userId := api.GetUserId(c)

	params := apiParams.GetParams(c)

	result, err := Service.getTransfers(userId, params)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")

	}
	return c.Status(200).JSON(result)
}

func UpdateTransfer(c *fiber.Ctx) error {
	req := new(Transfer)

	if err := c.BodyParser(req); err != nil {
		return c.Status(500).SendString("internal server error")
	}

	updatedFields, err := api.GetUpdatedField(c.BodyRaw())
	if err != nil {
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	req.ID = c.Params("transferId")

	req.UserId = api.GetUserId(c)

	if err := Service.updateTransfer(req, updatedFields); err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}
	return c.Status(200).JSON(req)
}
