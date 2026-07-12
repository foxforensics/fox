// Package version will be set upon release.
package version

var Number = "0.0.0"
var Commit = "dev"

// ID returns to shortened commit id.
//
//goland:noinspection GoBoolExpressions
func ID() string {
	if len(Commit) > 0 {
		return Commit[:min(len(Commit), 6)]
	}

	return "none"
}
