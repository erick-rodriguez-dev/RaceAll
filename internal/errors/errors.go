package errors

import (
	"errors"
	"fmt"
)

var (
	// Common connection errors
	ErrConnectionClosed  = errors.New("connection is closed")
	ErrConnectionFailed  = errors.New("connection failed")
	ErrConnectionTimeout = errors.New("connection timeout")
	ErrAlreadyConnected  = errors.New("already connected")
	ErrNotConnected      = errors.New("not connected")

	// Common protocol errors
	ErrProtocolNotInitialized = errors.New("protocol not initialized")
	ErrInvalidMessageType     = errors.New("invalid message type")
	ErrInvalidData            = errors.New("invalid data")
	ErrBufferTooSmall         = errors.New("buffer too small")
	ErrPartialWrite           = errors.New("partial write")
	ErrPartialRead            = errors.New("partial read")

	// Common timeout errors
	ErrTimeout      = errors.New("operation timeout")
	ErrReadTimeout  = errors.New("read timeout")
	ErrWriteTimeout = errors.New("write timeout")

	// Common validation errors
	ErrInvalidCarIndex       = errors.New("invalid car index")
	ErrInvalidDriverIndex    = errors.New("invalid driver index")
	ErrInvalidSessionType    = errors.New("invalid session type")
	ErrInvalidSessionPhase   = errors.New("invalid session phase")
	ErrInvalidCarLocation    = errors.New("invalid car location")
	ErrInvalidEventType      = errors.New("invalid event type")
	ErrInvalidDriverCategory = errors.New("invalid driver category")

	// Common shared memory errors
	ErrSharedMemoryNotFound = errors.New("shared memory not found")
	ErrSharedMemoryAccess   = errors.New("shared memory access denied")
	ErrSharedMemoryMap      = errors.New("failed to map shared memory")
	ErrSharedMemoryUnmap    = errors.New("failed to unmap shared memory")

	// Common I/O errors
	ErrReadFailed     = errors.New("read operation failed")
	ErrWriteFailed    = errors.New("write operation failed")
	ErrEncodingFailed = errors.New("encoding failed")
	ErrDecodingFailed = errors.New("decoding failed")
)

type AppError struct {
	Module  string // Module where the error occurred (broadcast, shared-memory, etc.)
	Op      string // Operation that failed
	Err     error  // Underlying error
	Context string // Optional additional context
}

func (e *AppError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%s.%s: %s: %v", e.Module, e.Op, e.Context, e.Err)
	}
	return fmt.Sprintf("%s.%s: %v", e.Module, e.Op, e.Err)
}

// Unwrap allows errors.Is and errors.As to work correctly
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewError creates a new AppError
func NewError(module, op string, err error) error {
	if err == nil {
		return nil
	}
	return &AppError{
		Module: module,
		Op:     op,
		Err:    err,
	}
}

// NewErrorWithContext creates a new AppError with additional context
func NewErrorWithContext(module, op string, err error, context string) error {
	if err == nil {
		return nil
	}
	return &AppError{
		Module:  module,
		Op:      op,
		Err:     err,
		Context: context,
	}
}

// ValidationError represents a validation error with specific details
type ValidationError struct {
	Module string // Module where the error occurred
	Field  string // Field that failed validation
	Value  any    // Value that caused the error
	Rule   string // Validation rule that failed
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: validation failed for field '%s': %s (value: %v)", e.Module, e.Field, e.Rule, e.Value)
}

func NewValidationError(module, field string, value any, rule string) error {
	return &ValidationError{
		Module: module,
		Field:  field,
		Value:  value,
		Rule:   rule,
	}
}

func WrapError(module, op string, err error) error {
	if err == nil {
		return nil
	}
	return NewError(module, op, err)
}

func WrapErrorf(module, op string, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	context := fmt.Sprintf(format, args...)
	return NewErrorWithContext(module, op, err, context)
}

func IsConnectionError(err error) bool {
	return errors.Is(err, ErrConnectionClosed) ||
		errors.Is(err, ErrConnectionFailed) ||
		errors.Is(err, ErrConnectionTimeout) ||
		errors.Is(err, ErrNotConnected)
}

func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

func IsTimeoutError(err error) bool {
	return errors.Is(err, ErrTimeout) ||
		errors.Is(err, ErrReadTimeout) ||
		errors.Is(err, ErrWriteTimeout) ||
		errors.Is(err, ErrConnectionTimeout)
}
