package commands

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"

	"ppfenning92/housekeeping/utils"
)

func isDirGitRepository() (bool, error) {
	cmd := exec.Command(
		"git",
		"-P",
		"status",
	)

	if err := cmd.Run(); err != nil {
		log.Debugf("Not a git repository")
		return false, err
	}

	return true, nil
}

func GetDefaultBranch() (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short")

	out, err := cmd.Output()
	if err != nil {
		log.Debug("Cannot determine default branch")
		return "", err
	}

	return strings.TrimSpace(
		strings.Replace(string(out), "origin/", "", 1),
	), nil
}

func sanitzeBranchName(name string, idx int, arr []string) string {
	nonBranchRegex := regexp.MustCompile(`[\*\s]`)

	return nonBranchRegex.ReplaceAllString(name, "")
}

func GetFeatureBranches() ([]string, error) {
	cmd := exec.Command("git", "-P", "branch", "-l")
	out, err := cmd.Output()
	if err != nil {
		log.Debug("Cannot get branches")
		return []string{}, err
	}

	defaultBranch, _ := GetDefaultBranch()

	allLocalBranches := utils.Map[string, string](
		strings.Split(string(out), "\n"),
		sanitzeBranchName,
	)

	featureBranches := utils.Filter(
		allLocalBranches,
		func(str string, idx int, arr []string) bool { return str != "" && str != defaultBranch },
	)
	return featureBranches, nil
}

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	out, err := cmd.Output()
	if err != nil {
		log.Debug("Cannot determine current branch")
		return "", nil
	}

	return strings.TrimSpace(string(out)), nil
}

func CheckoutGitBranch(branchName string) {
	cmd := exec.Command("git", "checkout", branchName)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Cannot checkout branch '%s'", branchName)
	}
}

func MergeGitBranch(branchName string, flags ...string) {
	cmd := exec.Command("git", "merge", branchName)
	_, err := cmd.Output()
	if err != nil {
		log.Warnf("Cannot merge branch '%s'. Aborting merge", branchName)
		abort := exec.Command("git", "merge", "--abort")
		abort.Run()
	}
}

func GitPush(branchName string) {
	cmd := exec.Command("git", "push", "--set-upstream", "origin", branchName)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Cannot push branch '%s'", branchName)
	}
}
