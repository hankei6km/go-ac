// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

package ac

import (
	"crypto/sha256"
	"io"
)

// runFuncType defines type of function that is used in funcOutput.
type runFuncType func(argv []string, outStream, errStream io.Writer) error

// FuncOutputBuilder adds properties to Output(Builder).
type FuncOutputBuilder interface {
	runFunc(runFuncType) OutputBuilder
}

// funcOutput implements Output by using the function(gocredits.Run()).
type funcOutput struct {
	baseOutput
	runFunc runFuncType
}

func (c *funcOutput) Flush() (hash []byte, err error) {
	modules, err := c.modules()
	if err != nil {
		return nil, wrapf(err, "erorr in ProgOutput.Flush")
	}
	_, err = c.writePruned(modules)
	if err != nil {
		return nil, wrapf(err, "erorr in ProgOutput.Flush")
	}
	h := sha256.New()
	w := io.MultiWriter(c.outStream, h)
	if err := c.runFunc([]string{c.workDir}, w, c.errStream); err != nil {
		return nil, wrapf(err, "erorr in ProgOutput.Flush - runFunc args(%s)", c.workDir)
	}
	return h.Sum(nil), nil
}

func newEmbedOutput(b *baseOutputBuilder) *funcOutput {
	return &funcOutput{
		baseOutput: *newBaseOutput(b),
		runFunc:    b.runFuncIntl,
	}
}
