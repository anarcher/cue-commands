package cmd

import (
	"io"
)

type Command interface {
	Stderr() io.Writer
	InOrStdin() io.Reader
	OutOrStdout() io.Writer
	OutOrStderr() io.Writer
}
