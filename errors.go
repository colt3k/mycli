package mycli

import (
	"fmt"
	"reflect"
)

type InvalidObjectError struct {
	Type reflect.Type
	Name string
}
func (e *InvalidObjectError) Error() string {
	if e.Type == nil {
		return "CLIFlag: nil"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "CLIflg: non-pointer "+e.Name+" " + e.Type.String() + ""
	}
	return "CLIFlag: nil " + e.Type.String() + ""
}

type InvalidValueError struct {
	Field string
	Value string
	Options interface{}
}
func (e *InvalidValueError) Error() string {
	if len(e.Value) == 0 {
		return fmt.Sprintf("Invalid value for '%s' VALUE: (empty)",e.Field)
	}

	return fmt.Sprintf("Invalid value for '%s' VALUE not valid '%s', VALID options are %v",e.Field,e.Value,e.Options)
}
