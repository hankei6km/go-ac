package ac

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Output provides functions to write CREADITS files from the binary file.
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

	CliOutput

	Branch() OutputBuilder
	Build() Output
}

type baseOutputBuilder struct {
	mu        *sync.Mutex
	goSumFile string
	workDir   string
	binary    string
	cli       string
	outStream io.Writer
	errStream io.Writer
}

func (b *baseOutputBuilder) GoSumFile(goSumFile string) OutputBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.goSumFile = goSumFile

	return b.branch()
}

func (b *baseOutputBuilder) WorkDir(workDir string) OutputBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.workDir = workDir

	return b.branch()
}

func (b *baseOutputBuilder) Binary(binary string) OutputBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.binary = binary

	return b.branch()
}

func (b *baseOutputBuilder) OutStream(outStream io.Writer) OutputBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.outStream = outStream

	return b.branch()
}

func (b *baseOutputBuilder) ErrStream(errStream io.Writer) OutputBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.errStream = errStream

	return b.branch()
}

func (b *baseOutputBuilder) Cli(cli string) OutputBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.cli = cli

	return b.branch()
}

func (b *baseOutputBuilder) branch() OutputBuilder {
	mu := b.mu
	// mu.Lock() 呼び出し元で lock されているときだけ実行するように注意.
	// defer mu.Unlock()

	b.mu = nil
	ret := *b // とりあえず
	ret.mu = &sync.Mutex{}

	b.mu = mu
	return &ret
}

func (b *baseOutputBuilder) Branch() OutputBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.branch()
}

func (b *baseOutputBuilder) Build() Output {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch {
	case b.cli != "":
		return newCliOutput(b)
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
		defer func() {
			errText := errStream.String()
			switch {
			case err != nil:
				w.CloseWithError(err)
			case errText != "":
				w.CloseWithError(errors.New(errText))
				// c.errStream
			}
			w.Close()
		}()
		cmd := exec.Command(c.modulesCmd, append(c.modulesArgs, c.binary)...) // 毎回appendはちょっともったいないか
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
		return nil, err
	case len(mods) == 0:
		return nil, errors.New("depndent not found")
	}
	return mods, nil
}

func (c *baseOutput) prune(modules []string) io.Reader {
	r, w := io.Pipe()

	go func() {
		var errClose error
		defer func() {
			w.CloseWithError(errClose)
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
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, c.prune(modules))
	if err != nil {
		return "", err
	}
	return outFile, err
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

		modulesCmd:  "go",
		modulesArgs: []string{"version", "-m"},

		builder: b.branch(),
	}
}

// NewOutputBuilder returns the instance of OutputBuilder.
func NewOutputBuilder() OutputBuilder {
	return &baseOutputBuilder{
		mu:        &sync.Mutex{},
		goSumFile: "go.sum",
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
}
