package cli

import "flag"

type Flag interface {
	Apply(set *flag.FlagSet)
	IsRequired() bool
}
