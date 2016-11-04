package preflight

import (
	"bytes"
	"os/exec"
	"syscall"
	"testing"
)

type PlaybookTest struct {
	// inputs
	Path string
	// expected outputs
	ExitCode int
	Output   []string // zero or more strings that should be in the output
}

func (p PlaybookTest) Run(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("ansible-playbook", p.Path)
	b, err := cmd.CombinedOutput()
	if p.ExitCode == 0 {
		if err != nil {
			t.Errorf("playbook execution failed: %v", err)
		}
	} else {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if gotCode := exitErr.Sys().(syscall.WaitStatus).ExitStatus(); gotCode != p.ExitCode {
				t.Errorf("got exit code %v, want %v", gotCode, p.ExitCode)
			}
		} else {
			t.Errorf("unexpected error: %#v", err)
		}
	}
	for _, s := range p.Output {
		if !bytes.Contains(b, []byte(s)) {
			t.Errorf("got:\n%s\nwant to contain %q", b, s)
		}
	}
}

// note: TestPing and TestFail below are just placeholders. The idea is to
// replace them with tests that call more intesting playbooks. However, the
// initial structure of the tests may be just like that: run a command, capture
// output, check for errors. I believe in most cases we will expect a non-nil
// error from running the command (non-zero exit code).

func TestPing(t *testing.T) {
	PlaybookTest{
		Path:   "test_ping.yml",
		Output: []string{"[test ping]"},
	}.Run(t)
}

func TestFail(t *testing.T) {
	PlaybookTest{
		Path:     "test_fail.yml",
		ExitCode: 1,
		Output:   []string{"[test fail]", `"msg": "Failed as requested from task"`},
	}.Run(t)
}
