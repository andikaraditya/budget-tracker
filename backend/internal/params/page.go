package params

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Paging struct {
	Limit  int `query:"limit,default=25"`
	Offset int `query:"offset,default=0"`
}

func getPaging(c *fiber.Ctx) *Paging {
	p := &Paging{
		Limit:  25,
		Offset: 0,
	}

	c.QueryParser(p)
	return p
}

func (p *Paging) Compose() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, p.Offset)
}
