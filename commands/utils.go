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
