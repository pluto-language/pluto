package evaluation

import (
	"fmt"
)

func Err(ctx *Context, msg, tag string, fmts ...interface{}) Object {
	panic(fmt.Sprintf("Errors cannot yet be thrown...\n%s - %s",
		fmt.Sprintf(msg, fmts...),
		tag,
	))
}

func IsErr(o Object) bool {
	return false
	// panic("Implement me :)")
}
