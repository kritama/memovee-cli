package version

import "fmt"

var (
	Version   = "dev"
	Revision  = "unknown"
	BuildTime = "unknown"
)

func String() string {
	return fmt.Sprintf(
		"memovee %s (revision %s, built %s)",
		Version,
		Revision,
		BuildTime,
	)
}
