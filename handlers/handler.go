package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/xframe-go/x/repository"
	"github.com/xframe-go/x/requests"
	"github.com/xframe-go/x/responses"
	"github.com/xframe-go/x/x"
	"gorm.io/gorm"
)

type Handler[M any, C repository.ToModel[M], U repository.ToUpdateModel, K comparable] struct {
	Repo repository.Interface[M, C, U, K]
	responses.Base
	opt *Option[M, C, U, K]
}

func NewHandler[M any, C repository.ToModel[M], U repository.ToUpdateModel, K comparable](
	repo repository.Interface[M, C, U, K],
	options ...WithOption[M, C, U, K],
) *Handler[M, C, U, K] {
	opt := defaultOption[M, C, U, K]()
	for _, option := range options {
		option(opt)
	}

	return &Handler[M, C, U, K]{
		Repo: repo,
		opt:  opt,
	}
}

type PaginationResp[M any] struct {
	Data  []M   `json:"data"`
	Total int64 `json:"total"`
}

func (h *Handler[M, C, U, K]) List(ctx echo.Context) error {
	var (
		c      = ctx.Request().Context()
		params = requests.ParseQueryParams(ctx)
	)

	list, total, err := h.Repo.List(c, params)
	if err != nil {
		return err
	}

	resp := PaginationResp[M]{
		Data:  list,
		Total: total,
	}

	return h.Success(ctx, resp)
}

func (h *Handler[M, C, U, K]) BatchList(ctx echo.Context) error {
	var (
		c      = ctx.Request().Context()
		params = requests.ParseQueryParams(ctx)
	)

	list, err := h.Repo.BatchList(c, params)
	if err != nil {
		return err
	}

	return h.Success(ctx, list)
}

func (h *Handler[M, C, U, K]) Create(c echo.Context) error {
	var (
		req = NewContext(c)
		ctx = c.Request().Context()
	)

	var create C
	if err := req.Validated(&create); err != nil {
		return h.Failed(c, err)
	}

	m := create.ToModel()

	err := x.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if h.opt.beforeCreate != nil {
			if err := h.opt.beforeCreate(req, tx, &create, &m); err != nil {
				return err
			}
		}

		if err := h.Repo.Create(tx, &m); err != nil {
			return err
		}

		if h.opt.afterCreated != nil {
			if err := h.opt.afterCreated(req, tx, &create, &m); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return h.Failed(c, err)
	}

	return h.Created(c, m)
}

func (h *Handler[M, C, U, K]) Show(ctx echo.Context) error {
	var (
		c      = ctx.Request().Context()
		params = requests.ParseQueryParams(ctx)

		id = h.opt.primaryKeyConverter(ctx.Param("id"))
	)

	show, err := h.Repo.Show(c, id, params)
	if err != nil {
		return h.Failed(ctx, err)
	}
	return h.Success(ctx, show)
}

func (h *Handler[M, C, U, K]) Update(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = NewContext(c)
		id  = h.opt.primaryKeyConverter(c.Param("id"))
	)

	var update U
	if err := req.Validated(&update); err != nil {
		return h.Failed(c, err)
	}

	err := x.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if h.opt.beforeUpdate != nil {
			if err := h.opt.beforeUpdate(req, tx, &update, id); err != nil {
				return err
			}
		}

		if err := h.Repo.Update(ctx, tx, id, update); err != nil {
			return err
		}

		if h.opt.afterUpdated != nil {
			if err := h.opt.afterUpdated(req, tx, &update, id); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return h.Failed(c, err)
	}

	return h.Created(c, update)
}

func (h *Handler[M, C, U, K]) Destroy(c echo.Context) error {
	var (
		req = NewContext(c)
		ctx = c.Request().Context()
		id  = strings.Split(c.Param("id"), ",")
		ids = make([]K, 0)
	)

	for _, s := range id {
		ids = append(ids, h.opt.primaryKeyConverter(s))
	}

	err := x.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if h.opt.beforeDestroy != nil {
			if err := h.opt.beforeDestroy(req, tx, ids); err != nil {
				return err
			}
		}

		if err := h.Repo.Destroy(ctx, tx, ids...); err != nil {
			return err
		}

		if h.opt.afterDestroyed != nil {
			if err := h.opt.afterDestroyed(req, tx, ids); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return h.Failed(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
