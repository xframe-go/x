package contracts

import "context"

type Request interface {
	context.Context

	Method() string
	Path() string
	Response() ResponseWriter
	Logger() Logger
	Bind(i interface{}) error
	FormValue(name string) string

	Param(name string) string
	ParamInt(name string) int
	ParamUint(name string) uint
	ParamInt32(name string) int32
	ParamUint32(name string) uint32
	ParamInt64(name string) int64
	ParamUint64(name string) uint64

	QueryParam(name string, def ...string) string
	QueryParamInt(name string, def ...int) int
	QueryParamUint(name string, def ...uint) uint
	QueryParamInt32(name string, def ...int32) int32
	QueryParamUint32(name string, def ...uint32) uint32
	QueryParamInt64(name string, def ...int64) int64
	QueryParamUint64(name string, def ...uint64) uint64
}
