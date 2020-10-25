package cmd

import (
	"testing"
)

func TestBuildTools(t *testing.T) {
	dir := "./testdata/script/hello/"
	tags := []string{}
	args := []string{}

	tools, err := BuildTools(dir, tags, args)
	if err != nil {
		t.Fatal(err)
	}
	if tools == nil {
		t.Fatal("inst is nil")
	}

	typ := "command"
	name := "print"

	//TODO(anarcher): testcases
	o := tools.Lookup(typ, name)
	if !o.Exists() {
		t.Fatalf("%s %s not found", typ, name)
	}

}
