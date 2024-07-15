package params

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Params struct {
	Filters []Filter
	Page    *Paging
	Sorts   *Sort
}

func GetParams(c *fiber.Ctx) *Params {
	return &Params{
		Filters: getFilters(c),
		Page:    getPaging(c),
		Sorts:   getSorts(c),
	}
}

func (f *Params) ComposeFilter(sb *strings.Builder, args []any) []any {
	if f == nil {
		return args
	}
	for _, filter := range f.Filters {
		args = append(args, filter.Value)
		switch strings.ToLower(filter.Operator) {
		case "eq":
			sb.WriteString(fmt.Sprintf(" AND %s = $%d ", filter.Column, len(args)))
		case "neq":
			sb.WriteString(fmt.Sprintf(" AND %s <> $%d ", filter.Column, len(args)))
		case "lt":
			sb.WriteString(fmt.Sprintf(" AND %s < $%d ", filter.Column, len(args)))
		case "gt":
			sb.WriteString(fmt.Sprintf(" AND %s > $%d ", filter.Column, len(args)))
		case "lte":
			sb.WriteString(fmt.Sprintf(" AND %s <= $%d ", filter.Column, len(args)))
		case "gte":
			sb.WriteString(fmt.Sprintf(" AND %s >= $%d ", filter.Column, len(args)))
		}
	}
	return args
}
