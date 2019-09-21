package ac

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Dist provide functions to write some CREDITS files for each binary files
// (ie. foo_linux_386, foo_darwin_amd64 -> CREDITS_linxu_386, CREDITS_darwin_amd64)
//
// Dist は各バイナリファイルから、それぞれ用の CREDITS ファイルを書き出す機能を提供する.
type Dist interface {
	Run() error
}

// DistBuilder builds Dist.
type DistBuilder interface {
	WorkDir(string) DistBuilder
	DistDir(string) DistBuilder
	OutDir(string) DistBuilder
	BaseName(string) DistBuilder
	Uniq(bool) DistBuilder

	OutputBuilder(OutputBuilder) DistBuilder

	OutStream(io.Writer) DistBuilder
	ErrStream(io.Writer) DistBuilder

	Branch() DistBuilder
	Build() Dist
}

type baseDistBuilder struct {
	mu       *sync.Mutex
	workDir  string
	distDir  string
	outDir   string
	baseName string
	uniq     bool

	outputBuilder OutputBuilder

	outStream io.Writer
	errStream io.Writer
}

func (b *baseDistBuilder) WorkDir(workDir string) DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.workDir = workDir

	return b.branch()
}

func (b *baseDistBuilder) DistDir(distDir string) DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.distDir = distDir

	return b.branch()
}

func (b *baseDistBuilder) OutDir(outDir string) DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.outDir = outDir

	return b.branch()
}

func (b *baseDistBuilder) BaseName(baseName string) DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.baseName = baseName

	return b.branch()
}

func (b *baseDistBuilder) Uniq(uniq bool) DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.uniq = uniq

	return b.branch()
}

func (b *baseDistBuilder) OutputBuilder(outputBuilder OutputBuilder) DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.outputBuilder = outputBuilder

	return b.branch()
}

func (b *baseDistBuilder) OutStream(outStream io.Writer) DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.outStream = outStream

	return b.branch()
}

func (b *baseDistBuilder) ErrStream(errStream io.Writer) DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.errStream = errStream

	return b.branch()
}

func (b *baseDistBuilder) branch() DistBuilder {
	mu := b.mu
	// mu.Lock() 呼び出し元で lock されているときだけ実行するように注意.
	// defer mu.Unlock()

	b.mu = nil
	ret := *b // とりあえず
	ret.mu = &sync.Mutex{}

	b.mu = mu
	return &ret
}

func (b *baseDistBuilder) Branch() DistBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.branch()
}

func (b *baseDistBuilder) Build() Dist {
	return newBaseDist(b)
}

type outputHash struct {
	outFileName string
	hash        []byte
}

type baseDist struct {
	workDir  string
	distDir  string
	outDir   string
	baseName string
	uniq     bool

	outputBuilder OutputBuilder

	outStream io.Writer
	errStream io.Writer

	builder DistBuilder

	// hash []outputHash
	hash []*outputHash
}

func (d *baseDist) output(distName string) error {
	s := DistSuffix(distName)
	// TODO: support replacement like as GoReleaser.
	pOs := s[0]
	pArch := s[1]

	outFileName := filepath.Join(d.outDir, d.baseName+"_"+pOs+"_"+pArch)
	out, err := os.Create(outFileName)
	if err != nil {
		return wrapf(err, "output creating the output file")
	}
	defer out.Close()

	o := d.outputBuilder.Binary(filepath.Join(d.distDir, distName)).OutStream(out).Build()
	hash, err := o.Flush()
	d.hash = append(d.hash, &outputHash{
		outFileName: outFileName,
		hash:        hash,
	})
	return err
}

func (d *baseDist) uniqByHash() (uniqed bool, err error) {
	l := len(d.hash)
	t := d.hash[0]
	p, _ := filepath.Split(t.outFileName)
	dstFileName := filepath.Join(p, d.baseName)
	for i := 1; i < l; i++ {
		if bytes.Equal(t.hash, d.hash[i].hash) == false {
			return false, nil
		}
	}
	for i := 1; i < l; i++ {
		if err := os.Remove(d.hash[i].outFileName); err != nil {
			return false, wrapf(err, "uniqByHash removeing files")
		}
	}
	if err := os.Rename(t.outFileName, dstFileName); err != nil {
		return false, wrapf(err, "uniqByHash renaming file")
	}
	return true, nil
}

func (d *baseDist) Run() error {
	dirs, err := ioutil.ReadDir(d.distDir)
	if err != nil {
		return wrapf(err, "Dist.Run")
	}
	for _, p := range dirs {
		if p.IsDir() {
			// go.sum を上書きしているので、並列で動かさないように注意.
			if err := d.output(p.Name()); err != nil {
				return wrapf(err, "Dist.Run")
			}
		}
	}
	if len(d.hash) == 0 {
		return wrapf(fmt.Errorf(" No %s file(s) has been created", d.baseName), "Dist.Run")
	}
	if d.uniq {
		_, err := d.uniqByHash()
		if err != nil {
			return wrapf(err, "Dist.Run")
		}
	}
	return nil
}

func newBaseDist(b *baseDistBuilder) *baseDist {
	return &baseDist{
		workDir:  b.workDir,
		distDir:  b.distDir,
		outDir:   b.outDir,
		baseName: b.baseName,
		uniq:     b.uniq,

		outputBuilder: b.outputBuilder.Branch().
			WorkDir(b.workDir),

		outStream: b.outStream,
		errStream: b.errStream,

		builder: b.branch(),

		hash: []*outputHash{},
	}
}

// NewDistBuilder returns the instance of DistBuilder.
func NewDistBuilder() DistBuilder {
	return &baseDistBuilder{
		mu:            &sync.Mutex{},
		baseName:      "CREDITS",
		uniq:          true,
		outputBuilder: NewOutputBuilder(),
		outStream:     os.Stdout,
		errStream:     os.Stderr,
	}
}
