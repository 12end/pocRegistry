package cli

import (
	"flag"
)

type IntFlag struct {
	Name     string
	Usage    string
	Value    int
	Required bool
}

func (f *IntFlag) Apply(set *flag.FlagSet) {
	set.Int(f.Name, f.Value, f.Usage)
}

func (f *IntFlag) IsRequired() bool {
	return f.Required
}
