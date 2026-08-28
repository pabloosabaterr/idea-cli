package main

import (
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	notesDir string

	lookingAt bool

	position int
	notePosition int

	folders []string
	notes []string

	notesCache map[string][]string
	previewCache map[string]string

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

	m.notesCache = make(map[string][]string)
	m.previewCache = make(map[string]string)

	m.loadNotes(0)

	return m
}

func (m model) Init() tea.Cmd {
	return tea.RequestWindowSize
}

