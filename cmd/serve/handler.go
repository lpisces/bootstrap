package serve

import (
	"github.com/labstack/echo"
)

type Handler interface {
	Serve(c echo.Context) (err error)
}
