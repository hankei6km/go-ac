package ac

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Output provides functions to write CREADITS files from the binary file.
//
// Output はバイナリファイルから CREADTIS ファイルを書き出す機能を提供する.
type Output interface {
	Write() error
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
}

func (c *baseOutput) modules() ([]string, error) {
	r, w := io.Pipe()
	go func() {
		cmd := exec.Command(c.modulesCmd, append(c.modulesArgs, c.binary)...) // 毎回appendはちょっともったいないか
		cmd.Stdout = w
		cmd.Stderr = c.errStream
		if err := cmd.Start(); err != nil {
			w.CloseWithError(err)
		}
		w.CloseWithError(cmd.Wait())
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
	return mods, scanner.Err()
}

func (c *baseOutput) prune(modules []string) error {
	// cmd:=os.
	return nil
}

func (c *baseOutput) Write() error {
	return nil
}

func newBaseOutput(b *baseOutputBuilder) *baseOutput {
	return &baseOutput{
		goSumFile:   b.goSumFile,
		workDir:     b.workDir,
		binary:      b.binary,
		outStream:   b.outStream,
		modulesCmd:  "go",
		modulesArgs: []string{"version", "-m"},
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
