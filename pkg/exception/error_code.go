// Package exception provides error code definition
package exception

// Common module (00) error codes definition.
var (
	ErrValidation     = 4000001
	ErrUnauthorized   = 4010001
	ErrInternalServer = 5000001
	ErrUnknownError   = 500
)

// Repository module (01) error codes definition.
var (
	ErrRepositoryNotFound     = 4040101
	ErrPublishToQueueFailed   = 5000201
	ErrPaginationParamInvalid = 4000201
)
