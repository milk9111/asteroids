package system

import "fmt"

type systemError struct {
	msg string
}

func newSystemError(msg string) systemError {
	return systemError{
		msg: fmt.Sprintf("system: %s", msg),
	}
}

func (err systemError) Error() string {
	return err.msg
}
