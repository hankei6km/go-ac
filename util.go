// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

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
			return wrapf(err, "ResetDir() removing dir")
		}
	}
	return os.Mkdir(name, perm)
}

// DistSuffix returns suffix  of d(ie. foo_bar_linxu_386 -> [linx 386])
func DistSuffix(d string) []string {
	s := strings.Split(d, "_")
	return s[len(s)-2:]
}

// ReplaceItem replaces s by r.
func ReplaceItem(r [][]string, s string) string {
	for _, r := range r {
		if r[0] == s {
			return r[1]
		}
	}
	return s
}
