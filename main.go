package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	glamour "charm.land/glamour/v2"
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

type editorFinishedMsg struct{
	path string
	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
	case editorFinishedMsg:
		delete(m.previewCache, msg.path)
		delete(m.notesCache, m.folders[m.position])

		pos := m.notePosition
		m.loadNotes(m.position)
		if pos < len(m.notes) {
			m.notePosition = pos
		}
    case tea.KeyPressMsg:
        switch msg.String() {
        case "ctrl+c":
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
					return editorFinishedMsg{path, err}
				})
			}
		case "q", "esc":
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

func (m model) renderPreview(path string) string {
	if !m.lookingAt {
		return "Nothing to preview ;)"
	}

	if preview, ok := m.previewCache[path]; ok {
		return preview
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "Nothing to preview ;)"
	}

	out, err := glamour.Render(string(b), "dark")
	if err != nil {
		out = string(b)
	}

	m.previewCache[path] = out
	return out
}

func (m model) View() tea.View {
	commandsStyle := lipgloss.NewStyle().
		Height(1).
		Width(m.width)

	command := commandsStyle.Render("")

	h := max(m.height-lipgloss.Height(command), 0)
	rest := max(m.width-20, 0)

	notesWidth := rest / 2
	previewWidth := rest - notesWidth

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
		Width(notesWidth)

	previewStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		Height(max(h-3, 0)).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(previewWidth)

	labelStyle := lipgloss.NewStyle().
		Bold(true)

	sidebar := sidebarStyle.Render(m.renderFolders())
	notes := notesStyle.Render(m.renderNotes())

	label := "  Preview:"
	credit := "Made with Love Pablo <3  "
	gap := max(previewWidth-lipgloss.Width(label)-lipgloss.Width(credit), 0)

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render(label),
		strings.Repeat(" ", gap),
		credit,
	)

	path := filepath.Join(m.notesDir, m.folders[m.position], m.notes[m.notePosition])

	previewPane := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		header,
		"",
		previewStyle.Render(m.renderPreview(path)),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, notes, previewPane)

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

