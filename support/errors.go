package support

type IllegalArgumentError struct {
	msg string
}

func (err IllegalArgumentError) Error() string {
	return err.msg
}

func NewIllegalArgumentError(msg string) *IllegalArgumentError {
	return &IllegalArgumentError{msg: msg}
}
