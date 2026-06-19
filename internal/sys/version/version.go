// Package version will be set upon release.
package version

var Number = "dev"
var Commit = "none"

// ID returns to shortened commit id.
func ID() string {
	return Commit[:min(len(Commit), 6)]
}
