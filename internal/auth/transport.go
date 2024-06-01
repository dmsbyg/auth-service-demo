package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHttpHandler(s Service) http.Handler {
	h := handler{s}
	e := echo.New()
	v1 := e.Group("/v1")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	authRoute := v1.Group("/auth")
	authRoute.POST("/register", h.HandleRegister)
	authRoute.POST("/login", h.HandleLogin)
	return e
}

type handler struct {
	service Service
}

func (h handler) HandleRegister(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, "register is not implemented yet!")
}

func (h handler) HandleLogin(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, "register is not implemented yet!")
}
