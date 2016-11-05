package preflight

import (
	"testing"
)

func TestUpgradeDependencyMissing(t *testing.T) {
	PlaybookTest{
		Path:     "check-yum-update/test-upgrade-dependency-missing.yml",
		ExitCode: 1,
		Output:   []string{"cannot yum update due to dependencies", "yum update will fail due to missing dependency"},
	}.Run(t)
}

func TestYumRepoBroken(t *testing.T) {
	PlaybookTest{
		Path:     "check-yum-update/test-yum-repo-broken.yml",
		ExitCode: 1,
		Output:   []string{"reports broken repo", "Error with yum repository configuration"},
	}.Run(t)
}

func TestYumRepoDisabled(t *testing.T) {
	PlaybookTest{
		Path:   "check-yum-update/test-yum-repo-disabled.yml",
		Output: []string{"nothing blocks a yum update"},
	}.Run(t)
}

func TestYumRepoUnreachable(t *testing.T) {
	PlaybookTest{
		Path:     "check-yum-update/test-yum-repo-unreachable.yml",
		ExitCode: 1,
		Output:   []string{"repo cannot reach its url", "Error getting data from at least one yum repository"},
	}.Run(t)
}
