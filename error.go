package tex

func wrapErr(err error, message string) error {
	return &ErrWithMessage{
		Cause:   err,
		Message: message,
	}
}

type ErrWithMessage struct {
	Cause   error
	Message string
}

func (e *ErrWithMessage) Error() string {
	return e.Message + ": " + e.Cause.Error()
}

func (e ErrWithMessage) Unwrap() error {
	return e.Cause
}

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
