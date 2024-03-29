// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

package ac

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Songmu/gocredits"
)

// Output provides functions to write the CREADITS file from the binary file.
//
// Output はバイナリファイルから CREADTIS ファイルを書き出す機能を提供する.
type Output interface {
	Flush() (hash []byte, err error)
}

// OutputBuilder builds CreaditsFile.
//
// 今回はそれほどコスト気にする必要はないので、各メソッドで深いコピー(ぽいこと)を行う.
type OutputBuilder interface {
	GoSumFile(string) OutputBuilder
	WorkDir(string) OutputBuilder
	Binary(string) OutputBuilder
	OutStream(io.Writer) OutputBuilder
	ErrStream(io.Writer) OutputBuilder

	ProgOutput
	FuncOutputBuilder

	Branch() OutputBuilder
	Build() Output
}

type baseOutputBuilder struct {
	goSumFile   string
	workDir     string
	binary      string
	prog        string
	runFuncIntl runFuncType
	outStream   io.Writer
	errStream   io.Writer

	modulesCmd  string
	modulesArgs []string
}

func (b *baseOutputBuilder) GoSumFile(goSumFile string) OutputBuilder {
	bb := b.branch()
	bb.goSumFile = goSumFile
	return bb
}

func (b *baseOutputBuilder) WorkDir(workDir string) OutputBuilder {
	bb := b.branch()
	bb.workDir = workDir
	return bb
}

func (b *baseOutputBuilder) Binary(binary string) OutputBuilder {
	bb := b.branch()
	bb.binary = binary
	return bb
}

func (b *baseOutputBuilder) OutStream(outStream io.Writer) OutputBuilder {
	bb := b.branch()
	bb.outStream = outStream
	return bb
}

func (b *baseOutputBuilder) ErrStream(errStream io.Writer) OutputBuilder {
	bb := b.branch()
	bb.errStream = errStream
	return bb
}

func (b *baseOutputBuilder) Prog(prog string) OutputBuilder {
	bb := b.branch()
	bb.prog = prog
	return bb
}

func (b *baseOutputBuilder) runFunc(runFunc runFuncType) OutputBuilder {
	bb := b.branch()
	bb.runFuncIntl = runFunc
	return bb
}

func (b *baseOutputBuilder) branch() *baseOutputBuilder {
	return &(*b) // とりあえず
}

func (b *baseOutputBuilder) Branch() OutputBuilder {
	return b.branch()
}

func (b *baseOutputBuilder) Build() Output {
	switch {
	case b.prog != "":
		return newProgOutput(b)
	case b.runFuncIntl != nil:
		return newEmbedOutput(b)
	}
	return newBaseOutput(b)
}

type baseOutput struct {
	goSumFile string
	workDir   string
	binary    string
	outStream io.Writer
	errStream io.Writer

	modulesCmd  string
	modulesArgs []string

	builder OutputBuilder // 今回はおそらくつかわない.
}

func (c *baseOutput) modules() ([]string, error) {
	r, w := io.Pipe()
	go func() {
		var err error
		errStream := &strings.Builder{}
		args := append(c.modulesArgs, c.binary)
		defer func() {
			errText := errStream.String()
			switch {
			case err != nil:
				w.CloseWithError(wrapf(err, "execute args(%s, %x)", c.modulesCmd, args))
			case errText != "":
				w.CloseWithError(wrapf(fmt.Errorf("%s", errText), "execute: args(%s, %x)", c.modulesCmd, args))
				// io.Copy(c.errStream, strings.NewReader(errText))
			}
			w.Close()
		}()
		cmd := exec.Command(c.modulesCmd, args...)
		cmd.Stdout = w
		cmd.Stderr = errStream
		err = cmd.Start()
		if err != nil {
			return
		}
		err = cmd.Wait()
	}()
	mods := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "\t") {
			t := strings.SplitN(l, "\t", 4)
			if t[1] == "dep" {
				mods = append(mods, t[2])
			}
		}
	}
	err := scanner.Err()
	switch {
	case err != nil:
		return nil, wrapf(err, "modules()")
	case len(mods) == 0:
		return nil, wrapf(fmt.Errorf("depenet module not found in '%s'", c.binary), "modules()")
	}
	return mods, nil
}

func (c *baseOutput) prune(modules []string) io.Reader {
	r, w := io.Pipe()

	go func() {
		var errClose error
		defer func() {
			if errClose != nil {
				w.CloseWithError(wrapf(errClose, "prune go.sum"))
				return
			}
			w.Close()
		}()

		in, err := os.Open(c.goSumFile)
		if err != nil {
			errClose = err
			return
		}
		defer in.Close()

		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			l := scanner.Text()
			t := strings.SplitN(l, " ", 2)[0]
			for _, m := range modules {
				if m == t {
					fmt.Fprintln(w, l)
				}
			}
		}
		errClose = scanner.Err()
	}()

	return r
}

func (c *baseOutput) writePruned(modules []string) (outFile string, err error) {
	outFile = filepath.Join(c.workDir, "go.sum")
	out, err := os.Create(outFile)
	if err != nil {
		return "", wrapf(err, "creating pruned file")
	}
	defer out.Close()
	_, err = io.Copy(out, c.prune(modules))
	if err != nil {
		return "", wrapf(err, "writing pruned file")
	}
	return outFile, nil
}

func (c *baseOutput) Flush() (hash []byte, err error) {
	return
}

func newBaseOutput(b *baseOutputBuilder) *baseOutput {
	return &baseOutput{
		goSumFile: b.goSumFile,
		workDir:   b.workDir,
		binary:    b.binary,
		outStream: b.outStream,
		errStream: b.errStream,

		modulesCmd:  b.modulesCmd,
		modulesArgs: b.modulesArgs,

		builder: b.branch(),
	}
}

// NewOutputBuilder returns the instance of OutputBuilder.
func NewOutputBuilder() OutputBuilder {
	return &baseOutputBuilder{
		goSumFile:   "go.sum",
		runFuncIntl: gocredits.Run,
		outStream:   os.Stdout,
		errStream:   os.Stderr,

		modulesCmd:  "go",
		modulesArgs: []string{"version", "-m"},
	}
}
