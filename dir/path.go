package dir

import (
	"os"
	"path/filepath"

	hd "github.com/mitchellh/go-homedir"
)

// Home returns the user's home directory
func Home() (string, error) {
	return hd.Dir()
}

// GetPath returns the value of $PLUTO
func GetPath() (string, error) {
	path := os.Getenv("PLUTO")

	if len(path) == 0 {
		home, err := Home()
		if err != nil {
			return "", err
		}

		return filepath.Join(home, "pluto"), nil
	}

	return path, nil
}
