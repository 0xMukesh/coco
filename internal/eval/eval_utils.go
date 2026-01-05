package eval

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/object"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{
		Value: true,
	}
	FALSE = &object.Boolean{
		Value: false,
	}
)

func nativeBoolToObjectBool(value bool) object.Object {
	if value == true {
		return TRUE
	} else {
		return FALSE
	}
}

func isTruthy(value object.Object) bool {
	if boolean, ok := value.(*object.Boolean); ok {
		return boolean.Value
	}

	if _, ok := value.(*object.Null); ok {
		return false
	}

	return true
}

func newErrorObject(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}

	return false
}
