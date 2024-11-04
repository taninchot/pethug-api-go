package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func BindJSONWithValidation(ctx *gin.Context, obj interface{}) (bool, error) {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		// Create a slice to hold error messages
		var errorMessages []string

		// Check if the error is a validation error
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldErr := range validationErrors {
				errorMessages = append(errorMessages, fieldErr.Error())
			}
		} else {
			// If it's not a validation error, capture the general error
			errorMessages = append(errorMessages, err.Error())
		}

		// Return the JSON error response
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": errorMessages})
		return false, err // Indicates that binding failed
	}
	return true, nil // Indicates success
}
