package cmd

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestHelloTool(t *testing.T) {
	// t.Skip()

	checkError := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	const path = "./testdata/script/hello"
	_ = os.Chdir(path)

	args := []string{"print"}

	cmd := &cobra.Command{
		Use:          "root",
		SilenceUsage: true,
	}
	cmd.Flags().StringArrayP(string(flagInject), "t", nil,
		"set the value of a tagged field")

	rootCmd := &Command{Command: cmd, root: cmd, cmd: cmd}

	c, err := New(rootCmd, args)
	checkError(err)

	b := &bytes.Buffer{}
	c.SetOutput(b)

	err = c.Run(context.Background())
	checkError(err)

	t.Error(err, "\n", b.String())

}
