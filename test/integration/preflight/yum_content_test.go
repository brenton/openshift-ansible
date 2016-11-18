package preflight

import (
	"testing"
)

func TestInstallMissingRequired(t *testing.T) {
	PlaybookTest{
		Path:     "yum_content/test-install-missing-required.yml",
		ExitCode: 1,
		Output:   []string{"Cannot install all of the necessary packages"},
	}.Run(t)
}

func TestUpgradeDependencyMissing(t *testing.T) {
	PlaybookTest{
		Path:     "yum_content/test-upgrade-dependency-missing.yml",
		ExitCode: 1,
		Output:   []string{"Could not perform a yum update", "Errors from dependency resolution"},
	}.Run(t)
}

func TestYumRepoBroken(t *testing.T) {
	PlaybookTest{
		Path:     "yum_content/test-yum-repo-broken.yml",
		ExitCode: 1,
		Output:   []string{"reports broken repo", "Error with yum repository configuration"},
	}.Run(t)
}

func TestYumRepoDisabled(t *testing.T) {
	PlaybookTest{
		Path:   "yum_content/test-yum-repo-disabled.yml",
		Output: []string{"nothing blocks a yum update"},
	}.Run(t)
}

func TestYumRepoUnreachable(t *testing.T) {
	PlaybookTest{
		Path:     "yum_content/test-yum-repo-unreachable.yml",
		ExitCode: 1,
		Output:   []string{"repo cannot reach its url", "Error getting data from at least one yum repository"},
	}.Run(t)
}

func TestCorrectAosVersion(t *testing.T) {
	PlaybookTest{
		Path:     "yum_content/test-correct-aos-version.yml",
		ExitCode: 0,
		Output:   []string{"version 3.2 matched"},
	}.Run(t)
}

func TestIncorrectAosVersion(t *testing.T) {
	PlaybookTest{
		Path:     "yum_content/test-incorrect-aos-version.yml",
		ExitCode: 1,
		Output:   []string{"Not all of the required packages are available at requested version"},
	}.Run(t)
}

func TestMultipleAosVersion(t *testing.T) {
	PlaybookTest{
		Path:     "yum_content/test-multiple-aos-version.yml",
		ExitCode: 1,
		Output:   []string{"Multiple minor versions of these packages"},
	}.Run(t)
}
