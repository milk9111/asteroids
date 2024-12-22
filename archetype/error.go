package archetype

import "fmt"

type archetypeError struct {
	msg string
}

func newArchetypeError(msg string) archetypeError {
	return archetypeError{
		msg: fmt.Sprintf("archetype: %s", msg),
	}
}

func (err archetypeError) Error() string {
	return err.msg
}
