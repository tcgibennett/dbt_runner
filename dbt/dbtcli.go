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
	return dbtcli("compile", args...)
}

func (dbt *Dbt) Execute(args ...string) *messages.Response {
	return dbtcli("run", args...)
}

func dbtcli(command string, args ...string) *messages.Response {
	subcmd := []string{command}
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

	if strings.Contains("Error", out.String()) {
		return &messages.Response{
			Status:  "Error",
			Message: out.String(),
		}
	}

	return &messages.Response{
		Status:  "Success",
		Message: out.String(),
	}
}
