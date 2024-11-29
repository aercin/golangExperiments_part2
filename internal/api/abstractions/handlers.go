package abstractions

import (
	"github.com/labstack/echo/v4"
)

type Handlers interface {
	AddProductToStock(c echo.Context) error

	GetStock(c echo.Context) error
}
