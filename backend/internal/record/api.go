package record

import (
	"errors"
	"fmt"

	"github.com/andikaraditya/budget-tracker/backend/internal/api"
	apiParams "github.com/andikaraditya/budget-tracker/backend/internal/params"
	"github.com/gofiber/fiber/v2"
)

func CreateRecord(c *fiber.Ctx) error {
	req := new(Record)

	if err := c.BodyParser(req); err != nil {
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	req.UserId = api.GetUserId(c)

	validationError := api.ValidateRequest(req)
	if len(validationError) > 0 {
		return api.SendErrorResponse(c, 400, validationError)
	}

	err := Service.createRecord(req)
	if err != nil {
		return api.SendErrorResponse(c, 500, "internal server error")
	}
	return c.Status(201).JSON(req)
}

func GetRecord(c *fiber.Ctx) error {
	req := new(Record)

	req.ID = c.Params("recordId")

	req.UserId = api.GetUserId(c)

	if err := Service.getRecord(req); err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return api.SendErrorResponse(c, 404, "record not found")
		}
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	return c.Status(200).JSON(req)
}

func GetRecords(c *fiber.Ctx) error {
	userId := api.GetUserId(c)

	params := apiParams.GetParams(c)

	result, err := Service.getRecords(userId, params)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")

	}
	return c.Status(200).JSON(result)
}

func UpdateRecord(c *fiber.Ctx) error {
	req := new(Record)

	if err := c.BodyParser(req); err != nil {
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	updatedFields, err := api.GetUpdatedField(c.BodyRaw())
	if err != nil {
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	req.ID = c.Params("recordId")

	req.UserId = api.GetUserId(c)

	if err := Service.updateRecord(req, updatedFields); err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}
	return c.Status(200).JSON(req)

}

func GetSummary(c *fiber.Ctx) error {
	req := new(Summary)

	userId := api.GetUserId(c)

	params := apiParams.GetParams(c)

	if err := Service.getSummary(req, userId, params); err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	return c.Status(200).JSON(req)
}
