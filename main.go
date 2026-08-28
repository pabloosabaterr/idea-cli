package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type model struct {
	notesDir string

	lookingAt bool

	position int
	notePosition int

	folders []string
	notes []string

	height int
	width int
}

func initModel() model {
	home, err := os.UserHomeDir()
	if err != nil {
		return model{}
	}

	var m model
	m.notesDir = filepath.Join(home, "notes")

	os.MkdirAll(m.notesDir, 0o755)

	entries, _ := os.ReadDir(m.notesDir)

	for _, e := range entries {
		if e.IsDir() {
			m.folders = append(m.folders, e.Name())
		}
	}

	m.loadNotes(0)

	return m
}

func (m model) Init() tea.Cmd {
	return tea.RequestWindowSize
}

func (m *model) loadNotes (i int) {
	m.notes = nil
	firstDirPath := filepath.Join(m.notesDir, m.folders[i])
	notes, _ := os.ReadDir(firstDirPath)
	for _, e := range notes {
		if filepath.Ext(e.Name()) == ".md" {
			m.notes = append(m.notes, e.Name())
		}
	}
}

func (m *model) UpdateNotes(movement int) {
	if movement == 0 {
		return
	}

	if m.lookingAt {
		next := m.notePosition + movement
		if next < 0 || next >= len(m.notes) {
			return
		}
		m.notePosition = next
	} else {
		next := m.position + movement
		if next < 0 || next >= len(m.folders) {
			return
		}
		m.position = next
		m.loadNotes(m.position)
	}
}

func (m model) RenderFolders() string {
	lines := make([]string, len(m.folders))
	for i, f := range m.folders {
		if i == m.position {
			lines[i] = "> " + f
		} else {
			lines[i] = "  " + f
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) RenderNotes() string {
	lines := make([]string, len(m.notes))
	for i, f := range m.notes {
		if i == m.notePosition && m.lookingAt {
			lines[i] = "> " + f
		} else {
			lines[i] = "  " + f
		}
	}
	return strings.Join(lines, "\n")
}

type editorFinishedMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
		case "j":
			m.UpdateNotes(1)
		case "tab":
			/*
			 Make this work for the notes as well
			 */
			if m.lookingAt {
				break
			}

			if m.position == len(m.folders) - 1 {
				m.position = 0
				m.loadNotes(m.position)
			} else {
				m.UpdateNotes(1)
			}
		case "k":
			m.UpdateNotes(-1)
		case "shift+tab":
			/*
			Look comment for case "tab"
			*/
			if m.lookingAt {
				break
			}

			if m.position == 0 {
				m.position = len(m.folders) - 1
				m.loadNotes(m.position)
			} else {
				m.UpdateNotes(-1)
			}
		case "enter":
			if !m.lookingAt {
				m.lookingAt = true
			} else {

				path := filepath.Join(m.notesDir, m.folders[m.position], m.notes[m.notePosition])
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vi"
				}

				return m, tea.ExecProcess(exec.Command(editor, path), func(err error) tea.Msg {
					return editorFinishedMsg{err}
				})
			}
		case "esc":
			m.lookingAt = false
        }
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
    }
    return m, nil
}

func (m model) View() tea.View {
	sidebarStyle := lipgloss.NewStyle().
		Width(20).
		Align(lipgloss.Left).
		PaddingLeft(2).
		Height(m.height).
		BorderStyle(lipgloss.RoundedBorder())

	notesStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		Height(m.height).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(m.width - 20)

	sidebar := sidebarStyle.Render(m.RenderFolders())
	notes := notesStyle.Render(m.RenderNotes())

	v := tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, sidebar, notes))
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

