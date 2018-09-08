package version

// Version components
const (
	Maj = "0"
	Min = "3"
	Fix = "1"
)

var (
	// Must be a string because scripts like dist.sh read this file.
	Version = "0.3.1"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
