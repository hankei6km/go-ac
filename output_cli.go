// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

package ac

import (
	"crypto/sha256"
	"io"
	"os/exec"
)

// ProgOutput adds properties to Output(Builder).
type ProgOutput interface {
	Prog(string) OutputBuilder
}

// progOutput implements Output by using external programs(cli tools).
type progOutput struct {
	baseOutput
	prog string
}

func (c *progOutput) Flush() (hash []byte, err error) {
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
	cmd := exec.Command(c.prog, c.workDir)
	cmd.Stdout = w
	cmd.Stderr = c.errStream
	if err := cmd.Start(); err != nil {
		return nil, wrapf(err, "erorr in ProgOutput.Flush - start args(%s)", c.workDir)
	}
	if err := cmd.Wait(); err != nil {
		return nil, wrapf(err, "erorr in ProgOutput.Flush - wait args(%s)", c.workDir)
	}
	return h.Sum(nil), nil
}

func newProgOutput(b *baseOutputBuilder) *progOutput {
	return &progOutput{
		baseOutput: *newBaseOutput(b),
		prog:       b.prog,
	}
}
