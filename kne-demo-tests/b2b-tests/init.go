package tests

import (
	"keysight/athena/tests/pkg/api"
)

var opts *api.Opts

func init() {
	opts = api.NewOpts()
	// will print warning if config doesn't exist or is bad
	opts.PatchFromFile("opts.json")
	opts.BindFlags()
}
