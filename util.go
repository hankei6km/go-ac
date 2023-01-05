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

// DistSuffix returns suffix of d(ie. linux_386 -> [linux 386], linux_amd64_v1 -> [linux amd64_v1])
func DistSuffix(d string) []string {
	s := strings.Split(d, "_")
	if len(s) < 2 {
		return s
	}
	return []string{s[0], strings.Join(s[1:], "_")}
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
