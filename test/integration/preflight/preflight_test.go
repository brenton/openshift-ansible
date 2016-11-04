package preflight

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestPing1(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("ansible-playbook", "test_ping1.yml")
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("playbook execution failed: %v", err)
	}
	want := "test ping 1"
	if !bytes.Contains(b, []byte(want)) {
		t.Errorf("got:\n%s\nwant to contain %q", b, want)
	}
}

func TestPing2(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("ansible-playbook", "test_ping2.yml")
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("playbook execution failed: %v", err)
	}
	want := "test ping 2"
	if !bytes.Contains(b, []byte(want)) {
		t.Errorf("got:\n%s\nwant to contain %q", b, want)
	}
}
