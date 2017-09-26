package dir

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	errNoSources = errors.New("use: no sources found")
)

// LocateSources finds the source files specified.
// pkg is a glob, such as "std/io" or "std/*".
// The package is located relative to dir.
func LocateSources(dir, pkg string) ([]string, error) {
	base := filepath.Join(dir, pkg)

	files, err := filepath.Glob(base)
	if err != nil {
		return nil, err
	}

	var newFiles []string

	for _, file := range files {
		stat, err := os.Stat(file)
		if err != nil {
			return nil, err
		}

		if !stat.IsDir() {
			if strings.HasSuffix(stat.Name(), ".pluto") {
				newFiles = append(newFiles, file)
			}

			continue
		}

		pfile := filepath.Join(file, stat.Name()+".pluto")

		pstat, err := os.Stat(pfile)
		if err != nil {
			return nil, nil
		}

		if !pstat.IsDir() && strings.HasSuffix(pstat.Name(), ".pluto") {
			newFiles = append(newFiles, pfile)
		}
	}

	if len(newFiles) == 0 {
		return nil, errNoSources
	}

	return newFiles, nil
}

// LocateRootSources finds the source files specified,
// relative to $PLUTO/packages.
func LocateRootSources(pkg string) ([]string, error) {
	path, err := GetPath()
	if err != nil {
		return nil, err
	}

	pkgs := filepath.Join(path, "packages")

	return LocateSources(pkgs, pkg)
}

// LocateAnySources first tries to locate sources as absolute,
// but if none are found, looks in $PLUTO/packages.
// For example, pkg=std/io will find the standard IO package,
// and /Users/.../packages/io/src/* will find the source files of it.
func LocateAnySources(pkg string) ([]string, error) {
	sources, err := LocateSources("/", pkg)
	if err != nil && err != errNoSources {
		return nil, err
	}

	if len(sources) == 0 {
		sources, err = LocateRootSources(pkg)
		if err != nil && err != errNoSources {
			return nil, err
		}
	}

	return sources, nil
}
