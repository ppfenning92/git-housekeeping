package commands

import (
	"os/exec"

	"github.com/charmbracelet/log"
)

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
