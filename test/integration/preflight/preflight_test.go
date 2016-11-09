package preflight

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
)

var prebuildPlays = []string{
	"yum_update/build.yml",
}

// TestMain ensure that the necessary build steps are performed
// (in parallel) prior to running the tests.
func TestMain(m *testing.M) {
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(len(prebuildPlays))
	for _, file := range prebuildPlays {
		go func(file string) {
			runPlaybook(file)
			wg.Done()
		}(file)
	}
	wg.Wait()
	os.Exit(m.Run())
}

// runPlaybook runs a single playbook and exits if it fails
func runPlaybook(file string) {
	cmd := exec.Command("ansible-playbook", file)
	cmd.Env = append(os.Environ(), "ANSIBLE_FORCE_COLOR=1")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Pre-run build playbook failed: %s\n%s", file, output)
		os.Exit(1)
	}
}

// A PlaybookTest executes a given Ansible playbook and checks the exit code and
// output contents.
type PlaybookTest struct {
	// inputs
	Path string
	// expected outputs
	ExitCode int
	Output   []string // zero or more strings that should be in the output
}

// Run runs the PlaybookTest.
func (p PlaybookTest) Run(t *testing.T) {
	// A PlaybookTest is intended to be run in parallel with other tests.
	t.Parallel()

	cmd := exec.Command("ansible-playbook", p.Path)
	cmd.Env = append(os.Environ(), "ANSIBLE_FORCE_COLOR=1")
	bytesOut, err := cmd.CombinedOutput()

	// Check exit code.
	if (err == nil) && (p.ExitCode != 0) {
		p.checkExitCode(t, 0, p.ExitCode, cmd, bytesOut)
	}
	if (err != nil) && (p.ExitCode == 0) {
		got, ok := getExitCode(err)
		if !ok {
			p.logCmdAndOutput(t, cmd, bytesOut)
			t.Fatalf("unexpected error: %#v", err)
		}
		p.checkExitCode(t, got, p.ExitCode, cmd, bytesOut)
	}

	// Check output contents.
	for _, search := range p.Output {
		if !bytes.Contains(bytesOut, []byte(search)) {
			p.logCmdAndOutput(t, cmd, bytesOut)
			t.Errorf("wanted that to contain %q", search)
		}
	}
}

// getExitCode returns an exit code and true if the exit code could be taken
// from err, false otherwise.
// The implementation is GOOS-specific, and currently only supports Linux.
func getExitCode(err error) (int, bool) {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return -1, false
	}
	waitStatus, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		return -1, false
	}
	return waitStatus.ExitStatus(), true
}

// checkExitCode marks the test as failed when got is different than want.
func (p PlaybookTest) checkExitCode(t *testing.T, got, want int, cmd *exec.Cmd, output []byte) {
	if got == want {
		return
	}
	p.logCmdAndOutput(t, cmd, output)
	t.Fatalf("got exit code %v, wanted %v", got, want)
}

// logCmdAndOutput logs how to re-run a command and a summary of the output of
// its last execution for debugging.
func (p PlaybookTest) logCmdAndOutput(t *testing.T, cmd *exec.Cmd, output []byte) {
	const maxLines = 10
	lines := bytes.Split(output, []byte("\n"))
	if len(lines) > maxLines {
		lines = append([][]byte{[]byte("...")}, lines[len(lines)-maxLines:len(lines)]...)
	}
	output = bytes.Join(lines, []byte("\n"))
	dir, err := filepath.Abs(cmd.Dir)
	if err != nil {
		panic(err)
	}
	t.Logf("command: (cd %s && %s)\noutput:\n%s", dir, strings.Join(cmd.Args, " "), output)
}

func TestBYOCentOS7(t *testing.T) {
	PlaybookTest{
		Path:     "playbooks/byo_centos7.yml",
		ExitCode: 2,
		Output: []string{
			// TODO(rhcarvalho): update test playbook to go past this error.
			"The error was: KeyError: 'ansible_default_ipv4'",
		},
	}.Run(t)
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
		ExitCode: 2,
		Output:   []string{"[test fail]", `"msg": "Failed as requested from task"`},
	}.Run(t)
}
