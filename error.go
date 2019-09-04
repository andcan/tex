package tex

func wrapErr(err error, message string) error {
	return &ErrWithMessage{
		Cause:   err,
		Message: message,
	}
}

// ErrWithMessage wraps error with an additional message
type ErrWithMessage struct {
	Cause   error
	Message string
}

// Error returns message followed by error
func (e *ErrWithMessage) Error() string {
	return e.Message + ": " + e.Cause.Error()
}

func (e *ErrWithMessage) Is(target error) bool {
	_, ok := target.(*ErrWithMessage)
	return ok
}

func (e ErrWithMessage) Unwrap() error {
	return e.Cause
}

// ErrInvalidSource describes an invalid source directory for templates
type ErrInvalidSource struct {
	Source string
}

func (e *ErrInvalidSource) Error() string {
	return "invalid source: " + e.Source
}

func (e *ErrInvalidSource) Is(target error) bool {
	_, ok := target.(*ErrInvalidSource)
	return ok
}

// ErrInvalidDest describes an invalid destination directory
type ErrInvalidDest struct {
	Dest string
}

func (e *ErrInvalidDest) Error() string {
	return "invalid dest: " + e.Dest
}

func (e *ErrInvalidDest) Is(target error) bool {
	_, ok := target.(*ErrInvalidDest)
	return ok
}

// ErrNoNoMatch is returned when there are no files matching source
type ErrNoNoMatch struct {
	Source string
}

func (e *ErrNoNoMatch) Error() string {
	return "no match for: " + e.Source
}

func (e *ErrNoNoMatch) Is(target error) bool {
	_, ok := target.(*ErrNoNoMatch)
	return ok
}
