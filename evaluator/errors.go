package evaluator

import (
	"fmt"

	"github.com/Zac-Garby/pluto/object"
)

func err(ctx *object.Context, msg, tag string, fmts ...interface{}) object.Object {
	panic(fmt.Sprintf("Errors cannot yet be thrown...\n%s - %s",
		fmt.Sprintf(msg, fmts...),
		tag,
	))
}

func isErr(o object.Object) bool {
	panic("Implement me :)")
}
