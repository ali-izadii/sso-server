package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func badRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "bad_request",
		Message: message,
	})
}

func unauthorized(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error:   "unauthorized",
		Message: message,
	})
}

func forbidden(c echo.Context, message string) error {
	return c.JSON(http.StatusForbidden, ErrorResponse{
		Error:   "forbidden",
		Message: message,
	})
}

func notFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "not_found",
		Message: message,
	})
}

func conflict(c echo.Context, message string) error {
	return c.JSON(http.StatusConflict, ErrorResponse{
		Error:   "conflict",
		Message: message,
	})
}

func internalError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "internal_error",
		Message: message,
	})
}

func success(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
	})
}
