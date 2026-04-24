package requests

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/xframe-go/x/validate"
)

type Request struct {
}

func (Request) Validated(ctx echo.Context, pointer any) error {
	if err := (&echo.DefaultBinder{}).BindBody(ctx, &pointer); err != nil {
		return err
	}

	return validate.Validated(pointer)
}

func ParseQueryParams(ctx echo.Context) QueryParams {
	filters := parseFilters(ctx)

	page := cast.ToInt(ctx.QueryParam("page"))
	if page == 0 {
		page = 1
	}

	pageSize := cast.ToInt(ctx.QueryParam("pageSize"))
	if pageSize == 0 {
		pageSize = 20
	}

	qp := QueryParams{
		Keyword:  ctx.QueryParam("_keyword"),
		Filters:  filters,
		Page:     page,
		PageSize: pageSize,
	}

	preload := ctx.QueryParam("preload")
	if len(preload) > 0 {
		qp.Preload = strings.Split(preload, ",")
	}
	return qp
}

func parseFilters(ctx echo.Context) []Filter {
	var filters []Filter
	params := ctx.QueryParams()

	for key, values := range params {
		if !strings.HasPrefix(key, "filter[") {
			continue
		}

		key = strings.TrimPrefix(key, "filter[")
		key = strings.TrimSuffix(key, "]")

		parts := strings.Split(key, "][")
		if len(parts) != 2 {
			continue
		}

		filters = append(filters, Filter{
			Field:    parts[0],
			Operator: Operator(parts[1]),
			Value:    convertValue(values[0]),
		})
	}

	return filters
}
