package cmd

import (
	"bytes"
	"io"
	"testing"
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

	tools, err := BuildTools(dir, tags, args)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.Buffer{}
	cmd := NewMockCommand(&buf)
	typ := "command"
	name := "print"

	err = DoTasks(cmd, typ, name, tools)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	t.Logf("commond output: %s", output)
	if output != "WORLD!\n" {
		t.Errorf("have: %s want: %s", output, "WORLD!")
	}
}
