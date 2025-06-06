package errs

import "errors"

var (
	ErrCommentContent         = errors.New("comment content must be between 1 and 2000 characters")
	ErrCommentNotFound        = errors.New("comment not found")
	ErrPostNotFound           = errors.New("post not found")
	ErrInternalServerError    = errors.New("internal server error")
	ErrInvalidInput           = errors.New("invalid input provided")
	ErrIncorrectCommentLength = errors.New("incorrect comment lenth")
	ErrCommentsNotAllowed     = errors.New("comments not allowed")
	ErrParentCommentNotFound  = errors.New("parent comment not found")
)
