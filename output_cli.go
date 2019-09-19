package ac

// CliOutput adds properties to Output(Builder).
type CliOutput interface {
	Cli(string) OutputBuilder
}

// cliOutput implements by using cli tools.
type cliOutput struct {
	baseOutput
	cli string
}

func newCliOutput(b *baseOutputBuilder) *cliOutput {
	return &cliOutput{
		baseOutput: *newBaseOutput(b),
		cli:        b.cli,
	}
}
