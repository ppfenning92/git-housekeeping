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

func RebaseGitBranch(branchName string, flags ...string) {
	cmd := exec.Command("git", append([]string{"rebase"}, flags...)...)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Cannot execute '%s'. Error: %s", cmd, err)
	}
}

func GitPush(branchName string) {
	cmd := exec.Command("git", "push")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Cannot push branch '%s'", branchName)
	}
}
