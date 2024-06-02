package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dmsbyg/auth-service-demo/internal/common"
	"github.com/dmsbyg/auth-service-demo/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHTTPHandler(s Service) http.Handler {
	h := handler{s}
	e := echo.New()
	v1 := e.Group("/v1")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Validator = validator.New()

	authRoute := v1.Group("/auth")
	authRoute.POST("/health", h.HandleHealth)
	authRoute.POST("/register", h.HandleRegister)
	authRoute.POST("/login", h.HandleLogin)
	return e
}

type handler struct {
	service Service
}

func (handler) HandleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
}

func (h handler) HandleRegister(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "bad request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	token, err := h.service.Register(c.Request().Context(), req.Email, []byte(req.Password))
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
			return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, ErrorResponse{"internal server error"})
	}

	return c.JSON(http.StatusCreated, RegisterResponse{Token: token})
}

func (h handler) HandleLogin(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "bad request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	res, err := h.service.Login(c.Request().Context(), req.Email, []byte(req.Password))
	if err != nil {
		if errors.Is(err, common.ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "email or password is incorrect"})
		}
		if errors.As(err, &common.InternalAppError{}) {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		}

		return c.JSON(http.StatusUnprocessableEntity, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: res})
}
