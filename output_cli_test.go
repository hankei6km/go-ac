package ac

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cliOutput_Flush(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err, "check")
	testDir := filepath.Join(cwd, "testdata")
	binDir := filepath.Join(testDir, "binDir")
	binFile := filepath.Join(binDir, "my_cmd")
	cliFile := filepath.Join(testDir, "dummy.sh")
	workDir := filepath.Join(testDir, "work_flush")
	goSumDir := filepath.Join(testDir, "goSum")
	tests := []struct {
		name    string
		builder OutputBuilder
		want    string
		wantErr bool
	}{
		{
			name: "basic",
			builder: NewOutputBuilder().
				WorkDir(workDir).
				Binary(binFile).
				GoSumFile(filepath.Join(goSumDir, "go.sum")).
				Cli(cliFile),
			want: "test: " + workDir + "\n",
		}, {
			name: "binary not exists",
			builder: NewOutputBuilder().
				WorkDir(workDir).
				Binary(filepath.Join(binDir, "foo")).
				GoSumFile(filepath.Join(goSumDir, "go.sum")).
				Cli(cliFile),
			wantErr: true,
		}, {
			name: "command not found",
			builder: NewOutputBuilder().
				WorkDir(workDir).
				Binary(binFile).
				GoSumFile(filepath.Join(goSumDir, "go.sum")).
				Cli(filepath.Join(testDir, "foo")),
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

			got := &strings.Builder{}
			gotHash, err := tt.builder.OutStream(got).Build().Flush()
			if (err != nil) != tt.wantErr {
				t.Errorf("cliOutput.Flush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Equal(t,
					fmt.Sprintf("%x", sha256.Sum256([]byte(tt.want))),
					fmt.Sprintf("%x", (gotHash)),
					"cliOutput.Flush()",
				)
				assert.Equal(t, tt.want, got.String(), "cliOutput.Flush() outStream")
			}
		})
	}
}
