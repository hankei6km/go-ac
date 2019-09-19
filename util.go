package ac

import "os"

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
