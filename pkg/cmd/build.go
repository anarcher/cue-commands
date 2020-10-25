package cmd

import (
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"
)

func BuildTools(dir string, tags, args []string) (*cue.Instance, error) {
	cfg := &load.Config{
		Dir:   dir,
		Tags:  tags,
		Tools: true,
	}
	binst := load.Instances(args, cfg)
	if len(binst) == 0 {
		return nil, nil
	}

	included := map[string]bool{}

	ti := binst[0].Context().NewInstance(binst[0].Root, nil)
	for _, inst := range binst {
		k := 0
		for _, f := range inst.Files {
			if strings.HasSuffix(f.Filename, "_tool.cue") {
				if !included[f.Filename] {
					_ = ti.AddSyntax(f)
					included[f.Filename] = true
				}
				continue
			}
			inst.Files[k] = f
			k++
		}
		inst.Files = inst.Files[:k]
	}

	insts, err := buildToolInstances(binst)
	if err != nil {
		return nil, err
	}

	inst := insts[0]
	if len(insts) > 1 {
		inst = cue.Merge(insts...)
	}

	inst = inst.Build(ti)
	return inst, inst.Err
}

func buildToolInstances(binst []*build.Instance) ([]*cue.Instance, error) {
	instances := cue.Build(binst)
	for _, inst := range instances {
		if inst.Err != nil {
			return nil, inst.Err
		}
	}

	// TODO check errors after the fact in case of ignore.
	for _, inst := range instances {
		if err := inst.Value().Validate(); err != nil {
			return nil, err
		}
	}
	return instances, nil
}
