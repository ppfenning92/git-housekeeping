package rebase

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"ppfenning92/housekeeping/commands"
	"ppfenning92/housekeeping/utils"
)

type model struct {
	status        int
	err           error
	workingDir    string
	auto          bool
	branches      []string
	mainBranch    string
	currentBranch string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return ""
}

func (m model) isDirGitRepository() {
	gst := exec.Command(
		"git",
		"-P",
		"status",
	)

	if _, err := gst.Output(); err != nil {
		log.Fatal(
			fmt.Sprintf("'%v' is not a git repository", m.workingDir),
		)
		os.Exit(1)
	}
}

func sanitzeBranchName(name string, idx int, arr []string) string {
	nonBranchRegex := regexp.MustCompile(`[\*\s]`)

	return nonBranchRegex.ReplaceAllString(name, "")
}

func (m *model) getFeatureBranches() {
	getBranchesCmd := exec.Command("git", "-P", "branch", "-l")
	out, err := getBranchesCmd.Output()
	if err != nil {
		log.Fatal("Cannot get branches")
		os.Exit(1)
	}

	fmt.Println("test", m.mainBranch)
	allLocalBranches := utils.Map[string, string](
		strings.Split(string(out), "\n"),
		sanitzeBranchName,
	)

	m.branches = utils.Filter(
		allLocalBranches,
		func(str string, idx int, arr []string) bool { return str != "" && str != m.mainBranch },
	)
}

func (m *model) getDefaultBranch() {
	getMainBranchCmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short")

	out, err := getMainBranchCmd.Output()
	if err != nil {
		log.Fatal("Cannot determine main branch")
		os.Exit(1)
	}

	m.mainBranch = strings.TrimSpace(strings.Replace(string(out), "origin/", "", 1))
}

func (m *model) getCurrentBranch() {
	getCurrentBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	out, err := getCurrentBranchCmd.Output()
	if err != nil {
		log.Fatal("Cannot determine current branch")
	}

	m.currentBranch = strings.TrimSpace(string(out))
}

func (m model) updateMain() {
	if m.currentBranch != m.mainBranch {
		commands.CheckoutGitBranch(m.mainBranch)
	}

	gitPull := exec.Command("git", "pull")
	if _, err := gitPull.Output(); err != nil {
		log.Fatalf("Could not update '%s'", m.mainBranch)
		os.Exit(1)
	}
}

func (m model) rebase() {
	for _, branch := range m.branches {
		// gco branch
		commands.CheckoutGitBranch(branch)
		commands.MergeGitBranch(m.mainBranch)
		commands.GitPush(branch)

		log.Infof("rebased %s", branch)
		// rebase
	}
}

func Rebase(auto bool) {
	path, err := os.Getwd()
	if err != nil {
		log.Error("Cannot determine current working directory. Error: %v", err)
		os.Exit(1)
	}

	data := model{
		workingDir: path,
		auto:       auto,
	}

	data.isDirGitRepository()
	data.getDefaultBranch()
	data.getFeatureBranches()
	data.getCurrentBranch()
	fmt.Println(data)
	data.updateMain()

	data.rebase()

	p := tea.NewProgram(data)
	if _, err := p.Run(); err != nil {
		commands.CheckoutGitBranch(data.currentBranch)
		fmt.Printf("Something went wrong! %v", err)
		os.Exit(1)
	}
	commands.CheckoutGitBranch(data.currentBranch)
	os.Exit(0)
}
