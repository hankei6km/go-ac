package ac

import (
	"os"
	"strings"
)

// ResetDir calls RemoveAll (if name is exists) and Mkdir.
func ResetDir(name string, perm os.FileMode) error {
	if _, err := os.Stat(name); err == nil {
		err = os.RemoveAll(name)
		if err != nil {
			return err
		}
	}
	return os.Mkdir(name, perm)
}

// DistSuffix returns suffix  of d(ie. foo_bar_linxu_386 -> [linx 386])
func DistSuffix(d string) []string {
	s := strings.Split(d, "_")
	return s[len(s)-2:]
}
