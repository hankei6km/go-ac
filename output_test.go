// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

package ac

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OutputBuilder_Build(t *testing.T) {
	tests := []struct {
		name    string
		builder OutputBuilder
		want    Output
	}{
		{
			name:    "basic",
			builder: NewOutputBuilder(),
			want:    &baseOutput{},
		}, {
			name:    "prog",
			builder: NewOutputBuilder().Prog("foo"),
			want:    &progOutput{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.builder.Build()
			assert.IsType(t, tt.want, got, "OutputBuilder.Build()")
		})
	}
}

func Test_baseOutput_modules(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err, "check")
	binDir := filepath.Join(cwd, "testdata", "binDir")
	tests := []struct {
		name    string
		builder OutputBuilder
		want    []string
		wantErr bool
	}{
		{
			name:    "basic",
			builder: NewOutputBuilder().Binary(filepath.Join(binDir, "my_cmd")),
			want: []string{
				"gopkg.in/yaml.v2",
			},
		}, {
			name:    "not exists",
			builder: NewOutputBuilder().Binary(filepath.Join(binDir, "foo")),
			wantErr: true,
		}, {
			name:    "not binary",
			builder: NewOutputBuilder().Binary(filepath.Join(binDir, "test.txt")),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.builder.Build()
			got, err := c.(*baseOutput).modules()
			assert.Equal(t, tt.want, got, "baseOutput.modules()")
			if (err != nil) != tt.wantErr {
				t.Errorf("baseOutput.modules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_baseOutput_writePruned(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err, "check")
	testDir := filepath.Join(cwd, "testdata")
	workDir := filepath.Join(testDir, "work_pruned")
	goSumDir := filepath.Join(testDir, "goSum")
	type args struct {
		modules []string
	}
	tests := []struct {
		name        string
		builder     OutputBuilder
		args        args
		want        string
		wantOutFile string
		wantErr     bool
	}{
		{
			name:    "basic",
			builder: NewOutputBuilder().WorkDir(workDir).GoSumFile(filepath.Join(goSumDir, "go.sum")),
			args: args{
				modules: []string{
					"gopkg.in/yaml.v2",
				},
			},
			want: `gopkg.in/yaml.v2 v2.2.2 h1:ZCJp+EgiOT7lHqUV2J862kp8Qj64Jo6az82+3Td9dZw=
gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
`,
			wantOutFile: filepath.Join(workDir, "go.sum"),
		}, {
			name:    "go.sum not exists",
			builder: NewOutputBuilder().WorkDir(workDir).GoSumFile(filepath.Join(goSumDir, "foo")),
			args: args{
				modules: []string{
					"gopkg.in/yaml.v2",
				},
			},
			wantErr: true,
		}, {
			name:    "missing workDir",
			builder: NewOutputBuilder().WorkDir(filepath.Join(testDir, "foo")).GoSumFile(filepath.Join(goSumDir, "go.sum")),
			args: args{
				modules: []string{
					"gopkg.in/yaml.v2",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ResetDir(workDir, os.ModePerm)
			assert.Nil(t, err, "check")
			defer os.RemoveAll(workDir)

			gotOutFile, err := tt.builder.Build().(*baseOutput).writePruned(tt.args.modules)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseOutput.writePruned() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantOutFile, gotOutFile, "baseOutput.writePruned() outFile")

			if err == nil {
				got, err := ioutil.ReadFile(filepath.Join(workDir, "go.sum"))
				assert.Nil(t, err, "check")
				assert.Equal(t, tt.want, string(got), "pruned go.sum")
			}
		})
	}
}
