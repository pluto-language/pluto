package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		return nil, fmt.Errorf("use: no sources found at %s", pkg)
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

// LocateAnySources first tries to locate sources in dir,
// but if none are found, looks in $PLUTO/packages.
func LocateAnySources(dir, pkg string) ([]string, error) {
	inDir, err := LocateSources(dir, pkg)
	if err != nil {
		return nil, err
	}

	if len(inDir) > 0 {
		return inDir, nil
	}

	return LocateRootSources(pkg)
}
