package category

import (
	"errors"
	"fmt"

	"github.com/andikaraditya/budget-tracker/backend/internal/api"
	apiParams "github.com/andikaraditya/budget-tracker/backend/internal/params"
	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {
	req := new(Category)

	if err := c.BodyParser(req); err != nil {
		return c.Status(500).SendString("internal server error")
	}

	validationError := api.ValidateRequest(req)
	if len(validationError) > 0 {
		return api.SendErrorResponse(c, 400, validationError)
	}

	err := Service.createCategory(req)
	if err != nil {
		return c.Status(500).SendString("internal server error")
	}
	return c.Status(201).JSON(req)
}

func GetCategory(c *fiber.Ctx) error {
	req := new(Category)

	req.ID = c.Params("categoryId")

	if err := Service.getCategory(req); err != nil {
		if errors.Is(err, api.ErrNotFound) {
			return api.SendErrorResponse(c, 404, "category not found")
		}
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	return c.Status(200).JSON(req)
}

func GetCategories(c *fiber.Ctx) error {
	params := apiParams.GetParams(c)

	result, err := Service.getCategories(params)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")

	}
	return c.Status(200).JSON(result)
}

func UpdateCategory(c *fiber.Ctx) error {
	req := new(Category)

	if err := c.BodyParser(req); err != nil {
		return c.Status(500).SendString("internal server error")
	}

	updatedFields, err := api.GetUpdatedField(c.BodyRaw())
	if err != nil {
		return api.SendErrorResponse(c, 500, "internal server error")
	}

	req.ID = c.Params("categoryId")

	if err := Service.updateCategory(req, updatedFields); err != nil {
		fmt.Printf("ERROR: %v", err)
		return api.SendErrorResponse(c, 500, "internal server error")
	}
	return c.Status(200).JSON(req)

}
