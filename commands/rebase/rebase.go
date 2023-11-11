package rebase

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"ppfenning92/housekeeping/commands"
	"ppfenning92/housekeeping/utils"
)

type (
	branch struct {
		name string
	}
	model struct {
		workingDir    string
		interactive   bool
		branches      []string
		list          list.Model
		defaultBranch string
		currentBranch string
	}
)

func (b branch) Title() string       { return b.name }
func (b branch) Description() string { return b.name }
func (b branch) FilterValue() string { return b.name }

func initModel() model {
	path, err := os.Getwd()
	if err != nil {
		log.Error("Cannot determine current working directory. Error: %v", err)
		os.Exit(1)
	}

	currentBranch, _ := commands.GetCurrentBranch()
	defaultBranch, _ := commands.GetDefaultBranch()
	brancheNames, _ := commands.GetFeatureBranches()
	branches := utils.Map[string, list.Item](
		brancheNames,
		func(branchName string, _i int, _ []string) list.Item { return branch{name: branchName} },
	)

	return model{
		workingDir:    path,
		currentBranch: currentBranch,
		defaultBranch: defaultBranch,
		branches:      brancheNames,
		list:          list.New(branches, list.NewDefaultDelegate(), 0, 0),
	}
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// const listHeight = 14
//
// var (
// 	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
// 	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
// 	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
// 	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
// 	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
// 	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
// )

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		case tea.KeyEnter.String(), tea.KeySpace.String():
			i, ok := m.list.SelectedItem().(branch)
			n := string(i)
			return m, nil

		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func (m model) updateMain() {
	if m.currentBranch != m.defaultBranch {
		commands.CheckoutGitBranch(m.defaultBranch)
	}

	gitPull := exec.Command("git", "pull")
	if _, err := gitPull.Output(); err != nil {
		log.Fatalf("Could not update '%s'", m.defaultBranch)
		os.Exit(1)
	}
}

func (m model) rebase() {
	for _, branch := range m.branches {
		// gco branch
		commands.CheckoutGitBranch(branch)
		commands.MergeGitBranch(m.defaultBranch)
		commands.GitPush(branch)

		log.Infof("rebased %s", branch)
		// rebase
	}
}

func Rebase(interactive bool) {
	model := initModel()
	model.interactive = interactive
	// data.updateMain()

	// data.rebase()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong! %v", err)
		os.Exit(1)
	}
	commands.CheckoutGitBranch(model.currentBranch)
	os.Exit(0)
}
