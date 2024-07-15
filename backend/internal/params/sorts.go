package params

import (
	"fmt"
	"strings"

	"github.com/andikaraditya/budget-tracker/backend/internal/db"
	"github.com/gofiber/fiber/v2"
)

type Sort struct {
	Column string
	Asc    bool
}

func getSorts(c *fiber.Ctx) *Sort {
	fs := c.Query("sort")
	if len(fs) > 0 {
		s := strings.SplitN(fs, ":", 2)
		if len(s) == 2 {
			col := s[0]
			asc := strings.ToLower(s[1]) == "asc"
			return &Sort{Column: col, Asc: asc}
		}
	}
	return nil
}
func (s *Sort) Compose() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("ORDER BY %s %s ", s.Column, db.Order(s.Asc))
}
