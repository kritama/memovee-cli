package command

import (
	"flag"
	"io"
)

type globalOptions struct {
	JSON           bool
	NoColor        bool
	NonInteractive bool
	Yes            bool
	ConfigPath     string
	Help           bool
}

func parseGlobalOptions(args []string) (globalOptions, []string, *Problem) {
	var options globalOptions
	flags := flag.NewFlagSet("memovee", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(&options.JSON, "json", false, "render a deterministic JSON result")
	flags.BoolVar(&options.NoColor, "no-color", false, "disable ANSI color")
	flags.BoolVar(&options.NonInteractive, "non-interactive", false, "disable prompts")
	flags.BoolVar(&options.Yes, "yes", false, "confirm an authorized mutation")
	flags.StringVar(&options.ConfigPath, "config", "", "read configuration from path")
	flags.BoolVar(&options.Help, "help", false, "show help")
	flags.BoolVar(&options.Help, "h", false, "show help")

	if err := flags.Parse(args); err != nil {
		return options, nil, &Problem{
			Category: CategoryUsage,
			Message:  err.Error(),
			Next:     "run `memovee help` for usage",
		}
	}

	return options, flags.Args(), nil
}
