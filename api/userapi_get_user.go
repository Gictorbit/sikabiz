package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// GetUser is the handler function to retrieve a user by ID
func (us *UserService) GetUser(c echo.Context) error {
	// Parse user ID from the request path
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id required"})
	}

	// Retrieve the user from the database
	user, err := us.db.GetUserById(c.Request().Context(), id)
	if err != nil {
		us.logger.Error("failed to get user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
	}

	// Return the user as JSON
	return c.JSON(http.StatusOK, user)
}
