package evaluation

import "fmt"

func Err(ctx *Context, msg, tag string, fmts ...interface{}) Object {
	msg = fmt.Sprintf(msg, fmts...)

	e := &Instance{
		Base: ctx.Get("Error"),
	}

	if e.Base == nil {
		panic("Since the prelude isn't loaded, errors cannot be thrown!")
	}

	e.Data = map[string]Object{
		"tag": &String{Value: tag},
		"msg": &String{Value: msg},
	}

	return e
}

func IsErr(o Object) bool {
	if instance, ok := o.(*Instance); ok {
		return instance.Base.(*Class).Name == "Error"
	}

	return false
}
