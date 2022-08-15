package cli

import (
	"flag"
	"strconv"
)

type FlagSet struct {
	Set *flag.FlagSet
}

func NewFlagSet(name string) *FlagSet {
	return &FlagSet{Set: flag.NewFlagSet(name, flag.ExitOnError)}
}

func (s FlagSet) String(name string) string {
	return lookupString(name, s.Set)
}

func lookupString(name string, set *flag.FlagSet) string {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := f.Value.String(), error(nil)
		if err != nil {
			return ""
		}
		return parsed
	}
	return ""
}

func (s FlagSet) Bool(name string) bool {
	return lookupBool(name, s.Set)
}

func lookupBool(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return parsed
	}
	return false
}

func (s FlagSet) Int(name string) int {
	return lookupInt(name, s.Set)
}

func lookupInt(name string, set *flag.FlagSet) int {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseInt(f.Value.String(), 0, 64)
		if err != nil {
			return 0
		}
		return int(parsed)
	}
	return 0
}

func (s FlagSet) CommandName() string {
	return s.Set.Name()
}
