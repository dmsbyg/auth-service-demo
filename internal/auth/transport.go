package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dmsbyg/auth-service-demo/internal/common"
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
	var req RegisterRequest
	token, err := h.service.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		var duplicateErr common.DuplicateError
		if errors.As(err, &duplicateErr) {
			return c.JSON(
				http.StatusUnprocessableEntity,
				ErrorResponse{
					Error: fmt.Sprintf("this %s is already registered", duplicateErr.Entity),
				},
			)
		}
		if errors.As(err, &common.InternalAppError{}) {
			return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, ErrorResponse{"internal server error"})
	}

	return c.JSON(http.StatusCreated, RegisterResponse{Token: token})
}

func (h handler) HandleLogin(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, "register is not implemented yet!")
}
