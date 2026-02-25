package utils

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// Response is a standard structure for JSON responses.
type Response struct{}

// ErrorDetail represents a structured error detail
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// ErrorResponseData represents a structured error response
type ErrorResponseData struct {
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

// JSONResponse sends a standard JSON response.
func (r *Response) JSONResponse(c echo.Context, statusCode int, data interface{}, message string, err string) error {
	response := map[string]interface{}{
		"data":    data,
		"message": message,
		"error":   err,
	}
	return c.JSON(statusCode, response)
}

// JSONErrorResponse sends a structured error response
func (r *Response) JSONErrorResponse(c echo.Context, statusCode int, message string, errors []ErrorDetail) error {
	errorData := ErrorResponseData{
		Message: message,
		Errors:  errors,
	}

	response := map[string]interface{}{
		"data":    nil,
		"message": "",
		"error":   errorData,
	}
	return c.JSON(statusCode, response)
}

// ParseValidationError parses validation error string and converts to structured errors
func (r *Response) ParseValidationError(errString string) []ErrorDetail {
	errors := []ErrorDetail{}

	// Parse error like "Key: 'CreateNodeRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
	lines := strings.Split(errString, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract field name
		fieldStart := strings.Index(line, "'")
		fieldEnd := strings.LastIndex(line, "'")

		if fieldStart != -1 && fieldEnd != -1 && fieldEnd > fieldStart {
			fieldPath := line[fieldStart+1 : fieldEnd]
			parts := strings.Split(fieldPath, ".")
			fieldName := parts[len(parts)-1]

			// Extract validation tag
			var message string
			if strings.Contains(line, "required") {
				message = fieldName + " is required"
			} else if strings.Contains(line, "email") {
				message = fieldName + " must be a valid email"
			} else if strings.Contains(line, "min") {
				message = fieldName + " is too short"
			} else if strings.Contains(line, "max") {
				message = fieldName + " is too long"
			} else {
				message = fieldName + " is invalid"
			}

			errors = append(errors, ErrorDetail{
				Field:   strings.ToLower(fieldName),
				Message: message,
				Code:    "validation_error",
			})
		}
	}

	// If no errors parsed, return generic error
	if len(errors) == 0 {
		errors = append(errors, ErrorDetail{
			Message: errString,
			Code:    "validation_error",
		})
	}

	return errors
}

// ValidationErrorResponse is a helper for validation error responses
func (r *Response) ValidationErrorResponse(c echo.Context, errString string) error {
	errors := r.ParseValidationError(errString)
	return r.JSONErrorResponse(c, http.StatusBadRequest, "Validation failed", errors)
}

// SuccessResponse is a helper for success responses.
func (r *Response) SuccessResponse(c echo.Context, data interface{}, message string) error {
	return r.JSONResponse(c, http.StatusOK, data, message, "")
}

// ErrorResponse is a helper for error responses.
func (r *Response) ErrorResponse(c echo.Context, statusCode int, err string) error {
	return r.JSONResponse(c, statusCode, nil, "", err)
}

// CreatedResponse is a helper for responses with HTTP 201 Created.
func (r *Response) CreatedResponse(c echo.Context, data interface{}, message string) error {
	return r.JSONResponse(c, http.StatusCreated, data, message, "")
}

// NoContentResponse is a helper for HTTP 204 No Content responses.
func (r *Response) NoContentResponse(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// ForbiddenResponse is a helper for HTTP 403 Forbidden responses.
func (r *Response) ForbiddenResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusForbidden, nil, "", err)
}

// UnauthorizedResponse is a helper for HTTP 401 Unauthorized responses.
func (r *Response) UnauthorizedResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusUnauthorized, nil, "", err)
}

// BadRequestResponse is a helper for HTTP 400 Bad Request responses.
func (r *Response) BadRequestResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusBadRequest, nil, "", err)
}

// InternalServerErrorResponse is a helper for HTTP 500 Internal Server Error responses.
func (r *Response) InternalServerErrorResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusInternalServerError, nil, "", err)
}

// ConflictResponse is a helper for HTTP 409 Conflict responses.
func (r *Response) ConflictResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusConflict, nil, "", err)
}

// NotFoundResponse is a helper for HTTP 404 Not Found responses.
func (r *Response) NotFoundResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusNotFound, nil, "", err)
}

// CustomResponse is a helper for custom responses with flexible status codes.
func (r *Response) CustomResponse(c echo.Context, statusCode int, data interface{}, message string, err string) error {
	return r.JSONResponse(c, statusCode, data, message, err)
}
