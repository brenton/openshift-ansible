package preflight

import (
	"testing"
)

func TestInstallMissingRequired(t *testing.T) {
	PlaybookTest{
		Path:     "yum_update/test-install-missing-required.yml",
		ExitCode: 1,
		Output:   []string{"Cannot install all of the necessary packages"},
	}.Run(t)
}

func TestUpgradeDependencyMissing(t *testing.T) {
	PlaybookTest{
		Path:     "yum_update/test-upgrade-dependency-missing.yml",
		ExitCode: 1,
		Output:   []string{"Could not perform yum update", "Errors from resolution"},
	}.Run(t)
}

func TestYumRepoBroken(t *testing.T) {
	PlaybookTest{
		Path:     "yum_update/test-yum-repo-broken.yml",
		ExitCode: 1,
		Output:   []string{"reports broken repo", "Error with yum repository configuration"},
	}.Run(t)
}

func TestYumRepoDisabled(t *testing.T) {
	PlaybookTest{
		Path:   "yum_update/test-yum-repo-disabled.yml",
		Output: []string{"nothing blocks a yum update"},
	}.Run(t)
}

func TestYumRepoUnreachable(t *testing.T) {
	PlaybookTest{
		Path:     "yum_update/test-yum-repo-unreachable.yml",
		ExitCode: 1,
		Output:   []string{"repo cannot reach its url", "Error getting data from at least one yum repository"},
	}.Run(t)
}
