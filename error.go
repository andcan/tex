package tex

func wrapErr(err error, message string) error {
	return &ErrWithMessage{
		Cause:   err,
		Message: message,
	}
}

// Used to provide additional details
type ErrWithMessage struct {
	Cause   error
	Message string
}

// Returns message followed by error
func (e *ErrWithMessage) Error() string {
	return e.Message + ": " + e.Cause.Error()
}

// Implement xerrors.Is
func (e *ErrWithMessage) Is(target error) bool {
	_, ok := target.(*ErrWithMessage)
	return ok
}

// Unwraps this error returning it's cause
func (e ErrWithMessage) Unwrap() error {
	return e.Cause
}

// Describes an invalid source directory for templates
type ErrInvalidSource struct {
	Source string
}

// Implement error
func (e *ErrInvalidSource) Error() string {
	return "invalid source: " + e.Source
}

// Implement xerrors.Is
func (e *ErrInvalidSource) Is(target error) bool {
	_, ok := target.(*ErrInvalidSource)
	return ok
}

// Describes an invalid destination directory
type ErrInvalidDest struct {
	Dest string
}

// Implement error
func (e *ErrInvalidDest) Error() string {
	return "invalid dest: " + e.Dest
}

// Implement xerrors.Is
func (e *ErrInvalidDest) Is(target error) bool {
	_, ok := target.(*ErrInvalidDest)
	return ok
}

// Returned when there are no files matching source
type ErrNoNoMatch struct {
	Source string
}

// Implement error
func (e *ErrNoNoMatch) Error() string {
	return "no match for: " + e.Source
}

// Implement xerrors.Is
func (e *ErrNoNoMatch) Is(target error) bool {
	_, ok := target.(*ErrNoNoMatch)
	return ok
}
