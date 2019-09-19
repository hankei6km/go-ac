package ac

import (
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
