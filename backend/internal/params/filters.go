package params

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Filter struct {
	Column           string `json:"column"`
	ColumnOverridden bool   `json:"-"`
	Operator         string `json:"operator"`
	Value            string `json:"value"`
}

func getFilters(c *fiber.Ctx) []Filter {
	var res []Filter
	fs := c.Query("filter")
	if len(fs) > 0 {
		fsArray := strings.Split(fs, ",")
		for _, f := range fsArray {
			s := strings.SplitN(f, ":", 3)
			if len(s) == 3 {
				col := s[0]
				op := s[1]
				val := s[2]
				res = append(res, Filter{Column: col, Operator: op, Value: val})
			}
		}
	}
	return res
}
