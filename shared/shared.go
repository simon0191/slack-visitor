package shared

import (
	"fmt"
)

type Options map[string]interface{}

type Error struct {
	Id     string
	Parent error
	Info   map[string]interface{}
}

func NewError(id string, parent error, additionalInfo ...Options) error {

	info := map[string]interface{}{}
	if len(additionalInfo) > 0 {
		info = additionalInfo[0]
	}

	return &Error{id, parent, info}
}

func (err Error) Error() string {
	msg := err.Id

	for key, val := range err.Info {
		msg += fmt.Sprintf("\n\t%s: %v", key, val)
	}

	if err.Parent != nil {
		msg += "\n\tParent: " + err.Parent.Error()
	}

	return msg
}
