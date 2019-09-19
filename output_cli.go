package ac

import (
	"crypto/sha256"
	"io"
	"os/exec"
)

// CliOutput adds properties to Output(Builder).
type CliOutput interface {
	Cli(string) OutputBuilder
}

// cliOutput implements by using cli tools.
type cliOutput struct {
	baseOutput
	cli string
}

func (c *cliOutput) Flush() (hash []byte, err error) {
	modules, err := c.modules()
	if err != nil {
		return nil, err
	}
	if _, err = c.writePruned(modules); err != nil {
		return nil, err
	}
	h := sha256.New()
	w := io.MultiWriter(c.outStream, h)
	cmd := exec.Command(c.cli, c.workDir)
	cmd.Stdout = w
	cmd.Stderr = c.errStream
	if err := cmd.Start(); err != nil {
	}
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func newCliOutput(b *baseOutputBuilder) *cliOutput {
	return &cliOutput{
		baseOutput: *newBaseOutput(b),
		cli:        b.cli,
	}
}
