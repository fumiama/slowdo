package slowdo

import "errors"

var (
	ErrWaitTimeTooShort = errors.New("wait time too short")
)
