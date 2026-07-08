package appresult

import (
	"errors"
	appError "github.com/The-Ogulgozel/Banking-system/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func GetHTTPStatus(code string) int {
	switch code {
	case "SE-00400":
		return http.StatusBadRequest
	case "SE-00401":
		return http.StatusUnauthorized
	case "SE-00403":
		return http.StatusForbidden
	case "SE-00404":
		return http.StatusNotFound
	case "SE-00409":
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func RespondAppError(c *gin.Context, err error) {
	var appErr *appError.AppError
	if errors.As(err, &appErr) {
		c.JSON(GetHTTPStatus(appErr.Code), ErrorResponse{
			Message: appErr.Message,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Message: "internal server error",
	})
}
