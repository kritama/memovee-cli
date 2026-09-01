package version

import "fmt"

type Info struct {
	Version   string
	Revision  string
	BuildTime string
}

var (
	Version   = "dev"
	Revision  = "unknown"
	BuildTime = "unknown"
)

func Current() Info {
	return Info{
		Version:   Version,
		Revision:  Revision,
		BuildTime: BuildTime,
	}
}

func (i Info) String() string {
	return fmt.Sprintf(
		"memovee %s (revision %s, built %s)",
		i.Version,
		i.Revision,
		i.BuildTime,
	)
}

func String() string {
	return Current().String()
}
