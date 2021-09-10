package util

import (
	"os"

	"github.com/open-policy-agent/opa/loader"
)

type LoaderFilter struct {
	Ignore []string
}

func (f LoaderFilter) Apply(abspath string, info os.FileInfo, depth int) bool {
	for _, s := range f.Ignore {
		if loader.GlobExcludeName(s, 1)(abspath, info, depth) {
			return true
		}
	}
	return false
}
