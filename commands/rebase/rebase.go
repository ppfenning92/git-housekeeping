package rebase

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type model struct {
	status     int
	err        error
	workingDir string
	auto       bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Info(msg)
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

	gst := exec.Command(
		"git",
		"-P",
		"status",
	)

	if _, err := gst.Output(); err != nil {
		log.Fatal(
			fmt.Sprintf("'%v' is not a git repository", data.workingDir),
		)
	}

	p := tea.NewProgram(data)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Something went wrong! %v", err)
		os.Exit(1)
	}
}
