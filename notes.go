package main

import (
	"os"
	"strings"
	"path/filepath"
)

func (m *model) loadNotes(i int) {
	m.notes = nil
	m.notePosition = 0
	if i < 0 || i >= len(m.folders) {
		return
	}

	firstDirPath := filepath.Join(m.notesDir, m.folders[i])
	notes, _ := os.ReadDir(firstDirPath)
	for _, e := range notes {
		if filepath.Ext(e.Name()) == ".md" {
			m.notes = append(m.notes, e.Name())
		}
	}
}

func (m *model) updateNotes(movement int) {
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

func (m model) renderNotes() string {
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

