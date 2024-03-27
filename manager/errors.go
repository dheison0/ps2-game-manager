package manager

import "errors"

var (
	ErrCoverNotFound      = errors.New("cover not found")
	ErrCoverRequestFailed = errors.New("failed to make a request to Github")
	ErrNameAlreadyExists  = errors.New("a game with the same name already exists")
	ErrAlreadyInstalled   = errors.New("this game is already installed(based on game image)")
)
