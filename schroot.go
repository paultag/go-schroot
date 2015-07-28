/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package schroot

import (
	"fmt"
	"strings"

	"os/exec"
)

// {{{ Command helpers

func getOutputLine(cmd *exec.Cmd) (string, error) {
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(out), " \n\t\r"), nil
}

// }}}

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

// vim: foldmethod=marker
