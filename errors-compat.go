// +build !go1.13
// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

package ac

import (
	"github.com/pkg/errors"
)

func wrapf(err error, format string, a ...interface{}) error {
	return errors.Wrapf(err, format, a...)
}
