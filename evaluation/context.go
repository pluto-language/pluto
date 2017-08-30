package evaluation

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Zac-Garby/pluto/ast"

	"gopkg.in/yaml.v2"
)

type Context struct {
	Store     map[string]Object
	Functions []*Function
	Packages  map[string]Package

	Outer *Context
}

func NewContext() *Context {
	return &Context{
		Store:    make(map[string]Object),
		Packages: make(map[string]Package),
	}
}

func (c *Context) Enclose() *Context {
	return &Context{
		Outer:    c,
		Packages: make(map[string]Package),
	}
}

func (c *Context) EncloseWith(args map[string]Object) *Context {
	return &Context{
		Store:    args,
		Outer:    c,
		Packages: make(map[string]Package),
	}
}

func (c *Context) Get(key string) Object {
	if obj, ok := c.Store[key]; ok {
		return obj
	} else if c.Outer != nil {
		return c.Outer.Get(key)
	} else {
		return nil
	}
}

func (c *Context) Assign(key string, obj Object) {
	if c.Outer != nil {
		if v := c.Outer.Get(key); v != nil {
			c.Outer.Assign(key, obj)
			return
		}
	}

	c.Store[key] = obj
}

func (c *Context) Declare(key string, obj Object) {
	c.Store[key] = obj
}

func (c *Context) Import(name string) Object {
	var root string

	if r, exists := os.LookupEnv("PLUTO"); exists {
		root = r
	} else {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}

		root = filepath.Join(usr.HomeDir, "pluto")
	}

	var pkgFile *os.File

	path := filepath.Join(root, "libraries", name)

	// if the package can be found in $PLUTO/libraries
	if _, err := os.Stat(path); err != nil {
		return Err(c, "package '%s' not found in %s", "ImportError", name, filepath.Join(root, "libraries"))
	} else {
		metaPath := filepath.Join(path, fmt.Sprintf("%s.yaml", name))
		pkgFile, err = os.Open(metaPath)

		if err != nil {
			return Err(c, "'%s' not found in %s", "ImportError", name+".yaml", path)
		}
	}

	pkgReader := bufio.NewReader(pkgFile)

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(pkgReader); err != nil {
		panic(err)
	}

	pkgData := buf.String()

	pkg := &Package{
		Context: c,
		Used:    false,
	}

	yaml.Unmarshal([]byte(pkgData), &pkg.Meta)

	return O_NULL
}

func (c *Context) AddFunction(fn Object) {
	if _, ok := fn.(*Function); !ok {
		panic("Not a function!")
	}

	c.Functions = append(c.Functions, fn.(*Function))
}

func (c *Context) GetFunction(pattern []ast.Expression) interface{} {
	for _, fn := range c.Functions {
		if len(pattern) != len(fn.Pattern) {
			continue
		}

		matched := true

		for i, item := range pattern {
			fItem := fn.Pattern[i]

			if itemID, ok := item.(*ast.Identifier); ok {
				if fItemID, ok := fItem.(*ast.Identifier); ok {
					if itemID.Value != fItemID.Value {
						matched = false
					}
				}
			} else if _, ok := item.(*ast.Argument); !ok {
				matched = false
			} else if _, ok := fItem.(*ast.Parameter); !ok {
				matched = false
			}
		}

		if matched {
			return fn
		}
	}

	if c.Outer != nil {
		return c.Outer.GetFunction(pattern)
	}

	for _, fn := range GetBuiltins() {
		if len(pattern) != len(fn.Pattern) {
			continue
		}

		matched := true

		for i, item := range pattern {
			fItem := fn.Pattern[i]

			if itemID, ok := item.(*ast.Identifier); ok {
				if fItem[0] != '$' {
					if itemID.Value != fItem {
						matched = false
					}
				}
			} else if _, ok := item.(*ast.Argument); !ok {
				matched = false
			} else if !(fItem[0] == '$') {
				matched = false
			}
		}

		if matched {
			return fn
		}
	}

	return nil
}
