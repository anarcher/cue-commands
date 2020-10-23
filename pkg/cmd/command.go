package cmd

import (
	"context"
	"io"
	"strings"

	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/token"

	"github.com/spf13/cobra"
)

type Command struct {
	// The currently active command.
	*cobra.Command

	root *cobra.Command

	// Subcommands
	cmd *cobra.Command

	hasErr bool
}

type errWriter Command

func (w *errWriter) Write(b []byte) (int, error) {
	c := (*Command)(w)
	c.hasErr = true
	return c.Command.OutOrStderr().Write(b)
}

func (c *Command) Stderr() io.Writer {
	return (*errWriter)(c)
}

func (c *Command) SetOutput(w io.Writer) {
	c.root.SetOut(w)
}

func (c *Command) SetInput(r io.Reader) {
	c.root.SetIn(r)
}

var ErrPrintedError = errors.New("terminating because of errors")

func (c *Command) Run(ctx context.Context) (err error) {
	// Three categories of commands:
	// - normal
	// - user defined
	// - help
	// For the latter two, we need to use the default loading.
	defer recoverError(&err)

	if err := c.root.Execute(); err != nil {
		return err
	}
	if c.hasErr {
		return ErrPrintedError
	}
	return nil
}

func recoverError(err *error) {
	switch e := recover().(type) {
	case nil:
	case panicError:
		*err = e.Err
	default:
		panic(e)
	}
	// We use panic to escape, instead of os.Exit
}

type panicError struct {
	Err error
}

func exit() {
	panic(panicError{ErrPrintedError})
}

func New(c *Command, args []string) (cmd *Command, err error) {
	defer recoverError(&err)
	cmd = c

	rootCmd := cmd.root
	if len(args) == 0 {
		return cmd, nil
	}
	rootCmd.SetArgs(args)

	//TODO(anarcher): addSubcommands

	err = cmd.cmd.ParseFlags(args)
	if err != nil {
		return nil, err
	}

	tags, err := cmd.cmd.Flags().GetStringArray(string(flagInject))
	if err != nil {
		return nil, err
	}

	args = cmd.cmd.Flags().Args()
	rootCmd.SetArgs(args)

	if c, _, err := rootCmd.Find(args); err == nil && c != nil {
		//TODO(anarcher): Need to understand this behavior.
		if c.Name() == args[0] {
			return cmd, nil
		}
	}

	if !isCommandName(args[0]) {
		return cmd, nil
	}

	tools, err := buildTools(cmd, tags, args[1:])
	if err != nil {
		return cmd, err
	}

	_, err = addCustom(cmd, rootCmd, commandSection, args[0], tools)
	if err != nil {
		err = errors.Newf(token.NoPos,
			`%s %q is not defined
Ensure commands are defined in a "_tool.cue" file.
Run 'cue help' to show available commands.`,
			commandSection, args[0],
		)
		return cmd, err
	}
	return cmd, nil
}

func isCommandName(s string) bool {
	return !strings.Contains(s, `/\`) &&
		!strings.HasPrefix(s, ".") &&
		!strings.HasSuffix(s, ".cue")
}
