// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

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

func Test_baseDist_Run_With_Func(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err, "check")
	testDir := filepath.Join(cwd, "testdata")
	distDir := filepath.Join(testDir, "distDir")
	outDir := filepath.Join(testDir, "outDir")
	workDir := filepath.Join(testDir, "work_dist")
	goSumDir := filepath.Join(testDir, "goSum")
	tests := []struct {
		name        string
		builder     DistBuilder
		fakeRunFunc runFuncType
		wantFiles   []string
		wantErr     bool
	}{
		{
			name: "basic",
			builder: NewDistBuilder().
				DistDir(distDir).
				OutDir(outDir).
				WorkDir(workDir),
			fakeRunFunc: func(argv []string, outStream, errStream io.Writer) error {
				// constant output.
				_, err := io.Copy(outStream, strings.NewReader("test"))
				return err
			},
			wantFiles: []string{"CREDITS"},
		}, {
			name: "multiple",
			builder: NewDistBuilder().
				DistDir(distDir).
				OutDir(outDir).
				WorkDir(workDir),
			fakeRunFunc: func(argv []string, outStream, errStream io.Writer) error {
				// different output.
				f, err := os.Open("/dev/random")
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = io.CopyN(outStream, f, 2048*5)
				return err
			},
			wantFiles: []string{"CREDITS_linux_386", "CREDITS_linux_amd64", "CREDITS_linux_amd64_v1"},
		}, {
			name: "replace",
			builder: NewDistBuilder().
				DistDir(distDir).
				OutDir(outDir).
				ReplaceOs([][]string{
					[]string{"linux", "Linux"},
					[]string{"windows", "Windows"},
				}).
				ReplaceArch([][]string{
					[]string{"386", "i386"},
				}).
				WorkDir(workDir),
			fakeRunFunc: func(argv []string, outStream, errStream io.Writer) error {
				// different output.
				f, err := os.Open("/dev/random")
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = io.CopyN(outStream, f, 2048*5)
				return err
			},
			wantFiles: []string{"CREDITS_Linux_i386", "CREDITS_Linux_amd64", "CREDITS_Linux_amd64_v1"},
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

			d := tt.builder.
				OutputBuilder(
					NewOutputBuilder().
						GoSumFile(filepath.Join(goSumDir, "go.sum")).
						runFunc(tt.fakeRunFunc),
				).
				Build()
			err = d.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("baseDist.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			files, err := ioutil.ReadDir(outDir)
			assert.Nil(t, err, "check")
			gotFileNames := make([]string, len(files))
			for i, f := range files {
				gotFileNames[i] = f.Name()
			}
			assert.ElementsMatch(t, tt.wantFiles, gotFileNames, "files")
		})
	}
}
func Test_baseDist_Run_With_Prog(t *testing.T) {
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
				DistDir(distDir).
				OutDir(outDir).
				WorkDir(workDir),
			fakeScript: `#!/bin/sh
# different output.
dd if=/dev/random count=5 status=none
`,
			wantNumFOutput: 3,
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
				OutputBuilder(
					NewOutputBuilder().
						GoSumFile(filepath.Join(goSumDir, "go.sum")).
						Prog(progFile),
				).
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
