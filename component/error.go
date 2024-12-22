package component

import "fmt"

type componentError struct {
	msg string
}

func newComponentError(msg string) componentError {
	return componentError{
		msg: fmt.Sprintf("component: %s", msg),
	}
}

func (err componentError) Error() string {
	return err.msg
}
