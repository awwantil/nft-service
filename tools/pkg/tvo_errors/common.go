package tvoerrors

import (
	"errors"
	"fmt"
)

// messages
const (
	FailedGetRows = "Failed to get rows"
	FailedScanRow = "Failed to scan row"
)

var (
	// Common
	ErrNotFound      = errors.New("not found")              // ErrNotFound indicates that the requested item was not found
	ErrInsertFailed  = errors.New("insert failed")          // ErrInsertFailed indicates that an insert operation failed
	ErrUpdateFailed  = errors.New("update failed")          // ErrUpdateFailed indicates that an update operation failed
	ErrDeleteFailed  = errors.New("delete failed")          // ErrDeleteFailed indicates that a delete operation failed
	ErrNicknameTaken = errors.New("nickname already taken") // ErrNicknameTaken represents an error indicating that the nickname is already in use
	ErrEmptyList     = errors.New("empty list")
	ErrCastClaims    = errors.New("failed to cast token to TokenData")
	ErrEmptyMetadata = errors.New("empty metadata")

	ErrInvalidRequestData = errors.New("invalid request data")
	ErrServerError        = errors.New("something went wrong")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidJWT         = errors.New("invalid or expired JWT")
	ErrMissingJWT         = errors.New("missing or malformed JWT")
	ErrConflict           = errors.New("conflict")

	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidPhone = errors.New("invalid phone number")

	ErrFollowUser   = errors.New("follow User")
	ErrUnfollowUser = errors.New("unfollow User")

	ErrNoNetwork      = errors.New("no network found")
	ErrTxTimeAgo      = errors.New("tx time is too ago")
	ErrTxCheckTimeout = errors.New("tx check timeout")

	ErrNoResults          = errors.New("no result found")
	ErrObjectWasNotDelete = errors.New("object was not delete")
	ErrInvalidUrl         = errors.New("invalid YouTube URL or missing video ID parameter")
	ErrInvalidResizeParam = errors.New("invalid resize param. format - WxH")
	ErrInvalidSizes       = errors.New("invalid size value")

	ErrAlreadyLiked   = errors.New("already liked")
	ErrAlreadyUnliked = errors.New("already unliked")
)

// Wrap оборачивает ошибки для прокидывания наверх по стеку вызова
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
