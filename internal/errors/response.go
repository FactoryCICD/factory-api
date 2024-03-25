package errors

import (
	"net/http"

	"github.com/FactoryCICD/factory-api/pkg/log"
	"github.com/go-chi/render"
)

// Type is the type for error types
type Type int64

const (
	// InternalServerError is the error message for internal server errors
	InternalServerError Type = 500
	// BadRequestError is the error message for bad request errors
	BadRequestError Type = 400
	// NotFoundError is the error message for not found errors
	NotFoundError Type = 404
	// UnauthorizedError is the error message for unauthorized errors
	UnauthorizedError Type = 401
	// ForbiddenError is the error message for forbidden errors
	ForbiddenError Type = 403
)

// RequestError is the error type for request errors
type RequestError struct {
	Err    error
	Status Type
	Msg    string
	Logger log.Logger
}

func (e *RequestError) Error() string {
	// Log error if status is 500
	if e.Status == InternalServerError {
		e.Logger.Error(e.Msg, "error", e.Err)
	}

	return e.Msg
}

// NewRequestError creates a new request error
func NewRequestError(err error, status Type, msg string, logger log.Logger) *RequestError {
	return &RequestError{
		Err:    err,
		Status: status,
		Msg:    msg,
		Logger: logger,
	}
}

// HandleError handles the error
func HandleError(w http.ResponseWriter, r *http.Request, err *RequestError) {
	errObj := struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	}{Message: err.Error(), StatusCode: int(err.Status)}

	switch err.Status {
	case BadRequestError:
		w.WriteHeader(http.StatusBadRequest)
	case NotFoundError:
		w.WriteHeader(http.StatusNotFound)
	case UnauthorizedError:
		w.WriteHeader(http.StatusUnauthorized)
	case ForbiddenError:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		errObj.Message = "Internal Server Error"
	}

	render.JSON(w, r, errObj)
}
