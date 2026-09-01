package command

import (
	"fmt"
	"io"

	"github.com/kritama/memovee-cli/internal/output"
	"github.com/kritama/memovee-cli/internal/version"
)

type IO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type helpResult struct {
	Command string `json:"command"`
	Usage   string `json:"usage"`
}

type versionResult struct {
	Command   string `json:"command"`
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	BuildTime string `json:"build_time"`
}

const usage = `Usage: memovee [global options] <command>

Commands:
  version  Print build version information
  help     Show this help

Global options:
  --json             Render a deterministic JSON result
  --no-color         Disable ANSI color
  --non-interactive  Disable prompts
  --yes              Confirm an authorized mutation
  --config <path>    Read configuration from path
  -h, --help         Show this help
`

func Run(args []string, streams IO) int {
	options, commandArgs, problem := parseGlobalOptions(args)
	renderer := output.NewRenderer(streams.Stdout, streams.Stderr, options.JSON)
	if problem != nil {
		return renderProblem(renderer, *problem)
	}

	if options.Help {
		return renderHelp(renderer)
	}

	if len(commandArgs) == 0 {
		return renderProblem(renderer, Problem{
			Category: CategoryUsage,
			Message:  "a command is required",
			Next:     "run `memovee help` for usage",
		})
	}

	switch commandArgs[0] {
	case "help":
		if len(commandArgs) != 1 {
			return unexpectedArguments(renderer, "help")
		}
		return renderHelp(renderer)
	case "version":
		if len(commandArgs) != 1 {
			return unexpectedArguments(renderer, "version")
		}
		return renderVersion(renderer)
	default:
		return renderProblem(renderer, Problem{
			Category: CategoryUsage,
			Message:  fmt.Sprintf("unknown command %q", commandArgs[0]),
			Next:     "run `memovee help` for usage",
		})
	}
}

func renderHelp(renderer output.Renderer) int {
	err := renderer.Success(helpResult{
		Command: "help",
		Usage:   usage,
	}, usage)
	return renderedExitCode(err)
}

func renderVersion(renderer output.Renderer) int {
	info := version.Current()
	err := renderer.Success(versionResult{
		Command:   "version",
		Version:   info.Version,
		Revision:  info.Revision,
		BuildTime: info.BuildTime,
	}, info.String()+"\n")
	return renderedExitCode(err)
}

func unexpectedArguments(renderer output.Renderer, command string) int {
	return renderProblem(renderer, Problem{
		Category: CategoryUsage,
		Message:  fmt.Sprintf("command %q does not accept arguments", command),
		Next:     "run `memovee help` for usage",
	})
}

func renderProblem(renderer output.Renderer, problem Problem) int {
	exitCode := ExitCode(problem.Category)
	err := renderer.Failure(output.Failure{
		Category: string(problem.Category),
		ExitCode: exitCode,
		Message:  problem.Message,
		Next:     problem.Next,
	})
	if err != nil {
		return ExitInternal
	}

	return exitCode
}

func renderedExitCode(err error) int {
	if err != nil {
		return ExitInternal
	}

	return ExitSuccess
}
