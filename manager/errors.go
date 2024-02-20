package manager

import "errors"

var (
	ErrCoverNotFound      = errors.New("cover not found")
	ErrNoPermission       = errors.New("we don't have access to write in this path")
	ErrCoverRequestFailed = errors.New("failed to make a request to Github")
)
