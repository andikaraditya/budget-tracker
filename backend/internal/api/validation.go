package api

import (
	"github.com/go-playground/validator"
)

func ValidateRequest(req any) []string {
	validate := validator.New()
	var errs []string
	if err := validate.Struct(req); err != nil {
		validationsErrors := err.(validator.ValidationErrors)
		for _, validateError := range validationsErrors {
			errs = append(errs, "error field: "+validateError.Field()+" detail: "+validateError.Tag()+" "+validateError.Param())
		}
		return errs
	}
	return errs
}
