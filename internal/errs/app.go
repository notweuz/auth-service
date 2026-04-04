package errs

import "google.golang.org/grpc/codes"

type AppError struct {
	Status  codes.Code
	Message string
	Reason  string
}

func (e AppError) Error() string {
	return e.Message + ": " + e.Reason
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
