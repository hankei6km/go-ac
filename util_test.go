// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

package ac

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResetDir(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err, "check")
	testDir := filepath.Join(cwd, "testdata")
	workDir := filepath.Join(testDir, "work_reset")
	dummyDir := filepath.Join(workDir, "dummy")

	type args struct {
		name string
		perm os.FileMode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				name: filepath.Join(workDir, "test"),
				perm: os.ModePerm,
			},
		}, {
			name: "remove all",
			args: args{
				name: filepath.Join(dummyDir),
				perm: os.ModePerm,
			},
		}, {
			name: "not exitst",
			args: args{
				name: filepath.Join(workDir, "nest", "test"),
				perm: os.ModePerm,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := os.Stat(workDir); err == nil {
				os.RemoveAll(workDir)
			}
			err := os.Mkdir(workDir, os.ModePerm)
			assert.Nil(t, err, "check")
			defer os.RemoveAll(workDir)
			err = os.Mkdir(dummyDir, os.ModePerm)
			assert.Nil(t, err, "check")

			if err := ResetDir(tt.args.name, tt.args.perm); (err != nil) != tt.wantErr {
				t.Errorf("ResetDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				stat, err := os.Stat(tt.args.name)
				assert.Nil(t, err, "check")
				assert.True(t, stat.IsDir(), "mkdir ", tt.name)
			}
		})
	}
}

func TestDistSuffix(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "basic",
			args: args{
				d: "linux_386",
			},
			want: []string{"linux", "386"},
		},
		{
			name: "sufix v1",
			args: args{
				d: "linux_amd64_v1",
			},
			want: []string{"linux", "amd64_v1"},
		},
		{
			name: "sufix v1 with prefix",
			args: args{
				d: "foo_linux_amd64_v1",
			},
			want: []string{"linux", "amd64_v1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DistSuffix(tt.args.d)
			assert.Equal(t, tt.want, got, "DistSuffix()")
		})
	}
}

func TestReplaceItem(t *testing.T) {
	type args struct {
		r [][]string
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				r: [][]string{
					[]string{"linux", "Linux"},
					[]string{"windows", "Windows"},
				},
				s: "windows",
			},
			want: "Windows",
		}, {
			name: "no match",
			args: args{
				r: [][]string{
					[]string{"linux", "Linux"},
					[]string{"windows", "Windows"},
				},
				s: "aix",
			},
			want: "aix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceItem(tt.args.r, tt.args.s)
			assert.Equal(t, tt.want, got, "ReplaceItem()")
		})
	}
}
