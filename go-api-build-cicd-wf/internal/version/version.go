// Package version holds build metadata injected at compile time via
// -ldflags "-X github.com/example/go-api/internal/version.Version=..."
// The CI workflow sets these from the git tag/SHA/build date; local
// `go build` without ldflags falls back to the "dev" defaults below.
package version

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Info is the JSON-serializable shape returned by /version.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}
