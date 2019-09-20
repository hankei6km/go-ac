// +build !go1.13

package ac

import (
	"github.com/pkg/errors"
)

func wrapf(err error, format string, a ...interface{}) error {
	return errors.Wrapf(err, format, a...)
}
