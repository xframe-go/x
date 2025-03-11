package xhttp

import (
	"github.com/labstack/echo/v4"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/utils"
	"github.com/xframe-go/x/utils/singleton"
	"net/http"
	"strconv"
	"time"
)

func createRequestFromEchoContext(ctx echo.Context, cfg contracts.ServerConfig, app contracts.Application) contracts.Request {
	return NewRequest(ctx, cfg, app)
}

type Request struct {
	cfg contracts.ServerConfig
	app contracts.Application
	echo.Context
	response singleton.Singleton[*XResponse]
}

func NewRequest(ctx echo.Context, cfg contracts.ServerConfig, app contracts.Application) contracts.Request {
	return &Request{
		Context: ctx,
		cfg:     cfg,
		app:     app,
	}
}

/*  Context */

func (r *Request) Deadline() (deadline time.Time, ok bool) {
	return r.Context.Request().Context().Deadline()
}

func (r *Request) Done() <-chan struct{} {
	return r.Context.Request().Context().Done()
}

func (r *Request) Err() error {
	return r.Context.Request().Context().Err()
}

func (r *Request) Value(key any) any {
	return r.Context.Request().Context().Value(key)
}

func (r *Request) Method() string {
	return r.Context.Request().Method
}

func (r *Request) Path() string {
	return r.Context.Path()
}

func (r *Request) Response() contracts.ResponseWriter {
	return r.response.GetOrSet(func() *XResponse {
		return &XResponse{
			req:    r,
			cfg:    r.cfg,
			code:   http.StatusOK,
			header: http.Header{},
		}
	})
}

func (r *Request) Logger() contracts.Logger {
	return r.app.Log().Context(r.Context.Request().Context())
}

/********************* query param ***************************/

func (r *Request) QueryParam(name string, def ...string) string {
	var defaultVal string
	if len(def) > 0 {
		defaultVal = def[0]
	}

	value := r.Context.QueryParam(name)
	if len(value) == 0 {
		return defaultVal
	}

	return value
}

func (r *Request) QueryParamInt(name string, def ...int) int {
	var defaultVal int
	if len(def) > 0 {
		defaultVal = def[0]
	}

	v := r.Context.QueryParam(name)
	if len(v) == 0 {
		return defaultVal
	}

	value, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return value
}

func (r *Request) QueryParamUint(name string, def ...uint) uint {
	defaultValues := utils.SliceTo[uint, int](def, func(t uint) int {
		return int(t)
	})
	return uint(r.QueryParamInt(name, defaultValues...))
}

func (r *Request) QueryParamInt32(name string, def ...int32) int32 {
	defaultValues := utils.SliceTo[int32, int](def, func(t int32) int {
		return int(t)
	})
	return int32(r.QueryParamInt(name, defaultValues...))
}

func (r *Request) QueryParamUint32(name string, def ...uint32) uint32 {
	defaultValues := utils.SliceTo[uint32, int](def, func(t uint32) int {
		return int(t)
	})
	return uint32(r.QueryParamInt(name, defaultValues...))
}

func (r *Request) QueryParamInt64(name string, def ...int64) int64 {
	defaultValues := utils.SliceTo[int64, int](def, func(t int64) int {
		return int(t)
	})
	return int64(r.QueryParamInt(name, defaultValues...))
}

func (r *Request) QueryParamUint64(name string, def ...uint64) uint64 {
	defaultValues := utils.SliceTo[uint64, int](def, func(t uint64) int {
		return int(t)
	})
	return uint64(r.QueryParamInt(name, defaultValues...))
}

/********************* route param ***************************/

func (r *Request) Param(name string) string {
	return r.Context.QueryParam(name)
}

func (r *Request) ParamInt(name string) int {
	val := r.Context.Param(name)
	value, _ := strconv.Atoi(val)
	return value
}

func (r *Request) ParamUint(name string) uint {
	return uint(r.ParamInt(name))
}

func (r *Request) ParamInt32(name string) int32 {
	return int32(r.ParamInt(name))
}

func (r *Request) ParamUint32(name string) uint32 {
	return uint32(r.ParamInt(name))
}

func (r *Request) ParamInt64(name string) int64 {
	return int64(r.ParamInt(name))
}

func (r *Request) ParamUint64(name string) uint64 {
	return uint64(r.ParamInt(name))
}
