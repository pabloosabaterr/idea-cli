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

func (m model) renderFolders() string {
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

type editorFinishedMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
		case "j":
			m.updateNotes(1)
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
				m.updateNotes(1)
			}
		case "k":
			m.updateNotes(-1)
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
				m.updateNotes(-1)
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
			if !m.lookingAt {
				return m, tea.Quit
			}
			m.lookingAt = false
        }
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
    }
    return m, nil
}

func (m model) renderPreview(path string) (t string) {
	bytes, err := os.ReadFile(path)
	if err != nil || !m.lookingAt {
		return
	}

	t = string(bytes)
	return
}

func (m model) View() tea.View {
	commandsStyle := lipgloss.NewStyle().
		Height(1).
		Width(m.width)

	command := commandsStyle.Render("")

	h := max(m.height-lipgloss.Height(command), 0)
	rest := max(m.width-20, 0)

	sidebarStyle := lipgloss.NewStyle().
		Width(20).
		Align(lipgloss.Left).
		PaddingLeft(2).
		Height(h).
		BorderStyle(lipgloss.RoundedBorder())

	notesStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		Height(h).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(rest / 2)

	previewStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		Height(h).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(rest - rest/2)

	sidebar := sidebarStyle.Render(m.renderFolders())
	notes := notesStyle.Render(m.renderNotes())

	path := filepath.Join(m.notesDir, m.folders[m.position], m.notes[m.notePosition])
	preview := previewStyle.Render(m.renderPreview(path))

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, notes, preview)

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Left, body, command))
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

