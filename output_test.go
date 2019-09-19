package ac

import (
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
			name:    "cli",
			builder: NewOutputBuilder().Cli("foo"),
			want:    &cliOutput{},
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
