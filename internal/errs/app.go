package errs

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Status  codes.Code
	Message string
	Reason  string
}

func (e AppError) Error() string {
	return e.Message + ": " + e.Reason
}

func ToGRPC(err error) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return status.Error(appErr.Status, appErr.Error())
	}

	return status.Error(codes.Internal, "internal server error")
}

func NotFound(message, reason string) *AppError {
	return &AppError{
		Status:  codes.NotFound,
		Message: message,
		Reason:  reason,
	}
}

func InternalError(message, reason string) *AppError {
	return &AppError{
		Status:  codes.Internal,
		Message: message,
		Reason:  reason,
	}
}

func Conflict(message, reason string) *AppError {
	return &AppError{
		Status:  codes.AlreadyExists,
		Message: message,
		Reason:  reason,
	}
}

func Unauthorized(message, reason string) *AppError {
	return &AppError{
		Status:  codes.Unauthenticated,
		Message: message,
		Reason:  reason,
	}
}

func MissingPermissions(message, reason string) *AppError {
	return &AppError{
		Status:  codes.PermissionDenied,
		Message: message,
		Reason:  reason,
	}
}
