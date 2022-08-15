package cli

import (
	"github.com/12end/requests"
	"go.uber.org/zap"
	"net/url"
)

type Command struct {
	// The name of the command
	Name string
	// A short description of the usage of this command
	Usage string
	// List of flags to parse
	Flags []Flag

	Action func(context *Context) (vulnerable bool, traceInfo []requests.TraceInfo)
}

type Context struct {
	Target *url.URL
	Trace  bool
	Logger *zap.Logger
	Params *FlagSet
}
