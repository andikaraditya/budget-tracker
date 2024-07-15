package source

import (
	"errors"
	"fmt"

	"github.com/andikaraditya/budget-tracker/backend/internal/api"
	apiParams "github.com/andikaraditya/budget-tracker/backend/internal/params"
	"github.com/gofiber/fiber/v2"
)

func CreateSource(c *fiber.Ctx) error {
	req := new(Source)

	if err := c.BodyParser(req); err != nil {
		return c.Status(500).SendString("internal server error")
	}

	validationError := api.ValidateRequest(req)
	if len(validationError) > 0 {
		return api.SendErrorResponse(c, 400, validationError)
	}

	req.UserId = api.GetUserId(c)

	err := Service.createSource(req)
	if err != nil {
		return c.Status(500).SendString("internal server error")
	}
	return c.Status(201).JSON(req)
}

func GetSource(c *fiber.Ctx) error {
	req := new(Source)

	req.ID = c.Params("sourceId")

	userId := api.GetUserId(c)

	if err := Service.getSource(req, userId); err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return api.SendErrorResponse(c, 404, "source not found")
		}
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	return c.Status(200).JSON(req)
}

func GetSources(c *fiber.Ctx) error {
	userId := api.GetUserId(c)

	params := apiParams.GetParams(c)

	result, err := Service.getSources(userId, params)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")

	}
	return c.Status(200).JSON(result)
}

func UpdateSource(c *fiber.Ctx) error {
	req := new(Source)

	if err := c.BodyParser(req); err != nil {
		return c.Status(500).SendString("internal server error")
	}

	updatedFields, err := api.GetUpdatedField(c.BodyRaw())
	if err != nil {
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	req.ID = c.Params("sourceId")

	req.UserId = api.GetUserId(c)

	if err := Service.updateSource(req, updatedFields); err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}
	return c.Status(200).JSON(req)

}
