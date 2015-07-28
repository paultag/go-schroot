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
	err := exec.Command("schroot", "-e", "-c", schroot.session).Run()
	schroot.active = false
	return err
}

func (schroot *Schroot) Command(cmd string, args ...string) (*exec.Cmd, error) {
	/* The semantics of foo(bar, baz, quix...) are so fucked up here
	 * it's not even funny. Let's do an append and unpack that. Sorry
	 * this looks like garbage, but it's either explicitly passing the
	 * elements into the call site, or use a foo... with nothing prefixing
	 * it. Sooo. Yeah.
	 *    -- PRT */
	if !schroot.active {
		return nil, fmt.Errorf("Error: Schroot is not online")
	}
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
	), nil
}

func NewSchroot(name string) (*Schroot, error) {
	schroot := Schroot{
		name:   name,
		active: false,
	}
	out, err := getOutputLine(exec.Command("schroot", "-b", "-c", name))
	if err != nil {
		return nil, err
	}
	schroot.session = out
	schroot.active = true

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
