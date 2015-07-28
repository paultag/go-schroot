package schroot

import (
	"fmt"
	"strings"

	"os/exec"
)

func getOutputLine(cmd *exec.Cmd) (string, error) {
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(out), " \n\t\r"), nil
}

type Schroot struct {
	location string
	name     string
	session  string
	active   bool
}

func (schroot *Schroot) End() error {
	return exec.Command("schroot", "-e", "-c", schroot.session).Run()
}

func (schroot *Schroot) Run(cmd string, args ...string) *exec.Cmd {
	/* The semantics of foo(bar, baz, quix...) are so fucked up here
	 * it's not even funny. Let's do an append and unpack that. Sorry
	 * this looks like garbage, but it's either explicitly passing the
	 * elements into the call site, or use a foo... with nothing prefixing
	 * it. Sooo. Yeah.
	 *    -- PRT */
	return exec.Command(
		"schroot",
		append(
			[]string{
				"-r",
				"-c", schroot.session,
				"--directory", "/",
				"--",
				cmd,
			}, args...)...,
	)
}

func NewSchroot(name string) (*Schroot, error) {
	schroot := Schroot{}
	out, err := getOutputLine(exec.Command("schroot", "-b", "-c", name))
	if err != nil {
		return nil, err
	}
	schroot.session = out

	out, err = getOutputLine(exec.Command(
		"schroot", "--location",
		"-c", fmt.Sprintf("session:%s", schroot.session),
	))
	if err != nil {
		return nil, err
	}
	schroot.location = out

	return &schroot, nil
}
