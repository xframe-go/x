package server

import (
	"github.com/labstack/echo/v4"
	"github.com/vmihailenco/msgpack/v5"
)

type MsgpackSerializer struct{}

func (m *MsgpackSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	return msgpack.NewEncoder(c.Response()).Encode(i)
}

func (m *MsgpackSerializer) Deserialize(c echo.Context, i interface{}) error {
	return msgpack.NewDecoder(c.Request().Body).Decode(i)
}
