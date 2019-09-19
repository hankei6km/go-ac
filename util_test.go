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
