package ac

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_baseDist_Run(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err, "check")
	testDir := filepath.Join(cwd, "testdata")
	distDir := filepath.Join(testDir, "distDir")
	outDir := filepath.Join(testDir, "outDir")
	progDir := filepath.Join(testDir, "work_prog")
	progFile := filepath.Join(progDir, "fake.sh")
	workDir := filepath.Join(testDir, "work_dist")
	goSumDir := filepath.Join(testDir, "goSum")
	tests := []struct {
		name           string
		builder        DistBuilder
		fakeScript     string
		wantNumFOutput int
		wantErr        bool
	}{
		{
			name: "basic",
			builder: NewDistBuilder().
				GoSumFile(filepath.Join(goSumDir, "go.sum")).
				DistDir(distDir).
				OutDir(outDir).
				WorkDir(workDir),
			fakeScript: `#!/bin/sh
# constant output.
echo test
`,
			wantNumFOutput: 1,
		}, {
			name: "multiple",
			builder: NewDistBuilder().
				GoSumFile(filepath.Join(goSumDir, "go.sum")).
				DistDir(distDir).
				OutDir(outDir).
				WorkDir(workDir),
			fakeScript: `#!/bin/sh
# different output.
dd if=/dev/random count=5 status=none
`,
			wantNumFOutput: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ResetDir(workDir, os.ModePerm)
			assert.Nil(t, err, "check")
			defer os.RemoveAll(workDir)

			err = ResetDir(outDir, os.ModePerm)
			assert.Nil(t, err, "check")
			defer os.RemoveAll(outDir)

			err = ResetDir(progDir, os.ModePerm)
			assert.Nil(t, err, "check")
			defer os.RemoveAll(progDir)

			err = func() error {
				// ファイルを閉じないと実行できないので.
				f, err := os.Create(progFile)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close() // ファイルはディレクトリごと消されるので個別削除はしない.
				_, err = io.Copy(f, strings.NewReader(tt.fakeScript))
				return err
			}()
			assert.Nil(t, err, "check")
			os.Chmod(progFile, 0700)

			d := tt.builder.
				OutputBuilder(NewOutputBuilder().Prog(progFile)).
				Build()
			err = d.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("baseDist.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			files, err := ioutil.ReadDir(outDir)
			assert.Nil(t, err, "check")
			assert.Equal(t, len(files), tt.wantNumFOutput, "Num of output files")
		})
	}
}
