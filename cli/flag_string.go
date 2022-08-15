package cli

import (
	"flag"
)

type StringFlag struct {
	Name     string
	Usage    string
	Value    string
	Required bool
}

func (f *StringFlag) Apply(set *flag.FlagSet) {
	set.String(f.Name, f.Value, f.Usage)
}

func (f *StringFlag) IsRequired() bool {
	return f.Required
}
