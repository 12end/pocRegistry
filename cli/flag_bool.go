package cli

import "flag"

type BoolFlag struct {
	Name     string
	Usage    string
	Value    bool
	Required bool
}

func (f *BoolFlag) Apply(set *flag.FlagSet) {
	set.Bool(f.Name, f.Value, f.Usage)
}

func (f *BoolFlag) IsRequired() bool {
	return f.Required
}
