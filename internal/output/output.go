package output

import (
	"encoding/json"
	"fmt"
	"io"
)

const SchemaVersion = "1"

type Renderer struct {
	stdout io.Writer
	stderr io.Writer
	json   bool
}

type successEnvelope struct {
	SchemaVersion string `json:"schema_version"`
	OK            bool   `json:"ok"`
	Result        any    `json:"result"`
}

type Failure struct {
	Category string `json:"category"`
	ExitCode int    `json:"exit_code"`
	Message  string `json:"message"`
	Next     string `json:"next,omitempty"`
}

type failureEnvelope struct {
	SchemaVersion string  `json:"schema_version"`
	OK            bool    `json:"ok"`
	Error         Failure `json:"error"`
}

func NewRenderer(stdout, stderr io.Writer, jsonOutput bool) Renderer {
	return Renderer{
		stdout: stdout,
		stderr: stderr,
		json:   jsonOutput,
	}
}

func (r Renderer) Success(result any, human string) error {
	if r.json {
		return encodeJSON(r.stdout, successEnvelope{
			SchemaVersion: SchemaVersion,
			OK:            true,
			Result:        result,
		})
	}

	_, err := io.WriteString(r.stdout, human)
	return err
}

func (r Renderer) Failure(failure Failure) error {
	if r.json {
		return encodeJSON(r.stdout, failureEnvelope{
			SchemaVersion: SchemaVersion,
			OK:            false,
			Error:         failure,
		})
	}

	if _, err := fmt.Fprintf(
		r.stderr,
		"error (%s): %s\n",
		failure.Category,
		failure.Message,
	); err != nil {
		return err
	}
	if failure.Next == "" {
		return nil
	}

	_, err := fmt.Fprintf(r.stderr, "next: %s\n", failure.Next)
	return err
}

func encodeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}
