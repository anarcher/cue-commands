package cmd

import (
	"bytes"
	"io"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
)

type MockCommand struct {
	buf *bytes.Buffer
}

func NewMockCommand(buf *bytes.Buffer) *MockCommand {
	c := &MockCommand{
		buf: buf,
	}
	return c
}
func (c *MockCommand) Stderr() io.Writer {
	return c.buf
}
func (c *MockCommand) InOrStdin() io.Reader {
	return c.buf
}
func (c *MockCommand) OutOrStdout() io.Writer {
	return c.buf
}
func (c *MockCommand) OutOrStderr() io.Writer {
	return c.buf
}

func TestDoTasks(t *testing.T) {
	dir := "./testdata/script/hello"
	tags := []string{}
	args := []string{}

	c := &load.Config{
		Dir:   dir,
		Tools: true,
		Tags:  tags,
	}
	insts := cue.Build(load.Instances(args, c))
	for _, i := range insts {
		if err := i.Value().Validate(); err != nil {
			t.Error(err)
		}
	}

	tools := insts[0]
	if len(insts) > 1 {
		tools = cue.Merge(insts...)
	}

	buf := bytes.Buffer{}
	cmd := NewMockCommand(&buf)
	typ := "command"
	name := "print"

	err := DoTasks(cmd, typ, name, tools)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	t.Logf("commond output: %s", output)
	if output != "WORLD!\n" {
		t.Errorf("have: %s want: %s", output, "WORLD!")
	}
}
