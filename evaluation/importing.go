package evaluation

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

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
		Context: c.Enclose(),
		Used:    false,
	}

	yaml.Unmarshal([]byte(pkgData), &pkg.Meta)

	pkg.Sources = getSourceFiles(path, pkg.Meta.Sources)

	pkg.Context.Declare("__title", &String{Value: pkg.Meta.Title})
	pkg.Context.Declare("__description", &String{Value: pkg.Meta.Description})
	pkg.Context.Declare("__version", &String{Value: pkg.Meta.Version})

	return O_NULL
}

func getSourceFiles(path string, globs []string) []string {
	var sources []string

	for _, glob := range globs {
		abs := filepath.Join(path, glob)
		if matches, err := filepath.Glob(abs); err != nil {
			panic(err)
		} else {
			sources = append(sources, matches...)
		}
	}

	return sources
}
