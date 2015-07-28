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

// Schroot object encapsulation. Contains lightweight amounts of state, so
// it's advised that you pass a pointer to this around.
type Schroot struct {
	location string
	name     string
	session  string
	active   bool
}

// Close the schroot session, killing the open session from schroot. To
// see how many you've left open because you forgot to call this method,
// run `schroot -l --all-sessions`. To close them, run
// `schroot -e -c session:xxxx` (where session:xxxx is the name, as seen
// in schroot -l output)
//
// It's advised that you `defer` this whenever creating a new schroot.
func (schroot *Schroot) End() error {
	err := exec.Command("schroot", "-e", "-c", schroot.session).Run()
	schroot.active = false
	return err
}

// Create a os/exec.Cmd to run the command inside the schroot chroot. Keep
// in mind that we're not currently smart enough to deal with things like
// the enviorn in a native way, so don't get too cute.
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

// Create a new schroot based on the chroot name. To get an idea
// of what chroots exist on the system, run `schroot -l --all-chroots`.
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
