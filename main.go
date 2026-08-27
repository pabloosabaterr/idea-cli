package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	folders []string
}

func initModel() model {
	home, err := os.UserHomeDir()
	if err != nil {
		return model{}
	}

	var m model
	notesDir := filepath.Join(home, "notes")

	os.MkdirAll(notesDir, 0o755)

	entries, _ := os.ReadDir(notesDir)

	for _, e := range entries {
		if e.IsDir() {
			m.folders = append(m.folders, e.Name())
		}
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder

	for _, name := range m.folders {
		b.WriteString(name + "\n")
	}

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
	}
}

