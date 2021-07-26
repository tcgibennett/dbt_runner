package dbt

import (
	"bytes"
	"os/exec"
	"strings"

	"stitchdata.com/dbt/runner/messages"
)

type Dbt struct {
	Path string
}

func New(path string) *Dbt {
	return &Dbt{
		Path: path,
	}
}

func (dbt *Dbt) Compile(args ...string) *messages.Response {
	subcmd := []string{"compile"}
	args = append(subcmd, args...)
	cmd := exec.Command("dbt", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return &messages.Response{
			Status:  "Error",
			Message: err.Error(),
		}
	}

	if strings.Contains("ERROR", out.String()) {
		return &messages.Response{
			Status: "Error",
			Message: out.String(),
		}
	}

	return &messages.Response{
		Status: "Success",
		Message: out.String(),
	}
}
