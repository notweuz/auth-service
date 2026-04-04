package errs

import (
	"errors"
)

var (
	RecordNotFound      = errors.New("record not found")
	RecordAlreadyExists = errors.New("record already exists")
)
