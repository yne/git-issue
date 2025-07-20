package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	files         []string
	issues        list.Model
	editor        viewport.Model
	editorFocused bool
	listWith      int
	windowWidth   int
	windowHeight  int
}

type listItem struct {
	filename string
	obj      map[string]interface{}
}

func (i listItem) safe(key string, def string) string {
	if s, ok := i.obj[key].(string); ok && s != "" {
		return s
	}
	return def
}

func (i listItem) Title() string {
	return i.safe("title", "no title")
}

func (i listItem) Description() string {
	return fmt.Sprintf("%s %s %s", i.safe("priority", ""), i.safe("assignee", "(no assignee)"), i.safe("description", ""))
}
func (i listItem) FilterValue() string { return i.filename }

func main() {
	files, _ := filepath.Glob(filepath.Join("issues", "*.json"))
	items := make([]list.Item, len(files))
	for i, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var obj map[string]interface{}
		json.Unmarshal(data, &obj)
		items[i] = listItem{filename: file, obj: obj}
	}

	issues := list.New(items, list.NewDefaultDelegate(), 1, 1)
	issues.SetShowHelp(false)
	issues.SetShowTitle(false)
	issues.SetShowStatusBar(false)
	issues.SetShowPagination(false)
	issues.SetShowFilter(false)
	m := model{
		listWith: 30,
		files:    files,
		issues:   issues,
		editor:   viewport.New(0, 0),
	}
	m.editor.SetContent("Welcome to git-issue CLI\nSelect an issue using up/down arrows\nthen open/close with left/right or enter/escape")

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "right":
			if !m.editorFocused {
				selectedItem, ok := m.issues.SelectedItem().(listItem)
				if ok {
					m.editorFocused = true
					m.editor.SetContent(formatJSONFile(selectedItem.filename))
					return m, nil
				}
			}
		case "esc", "q", "left":
			if m.editorFocused {
				m.editorFocused = false
				m.editor.SetContent("")
				return m, nil
			} else {
				return m, tea.Quit
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if !m.editorFocused {
		m.issues, cmd = m.issues.Update(msg)
	} else {
		m.editor, cmd = m.editor.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	m.editor.Width = m.windowWidth - m.listWith
	m.editor.Height = m.windowHeight
	m.issues.SetWidth(m.listWith)
	m.issues.SetHeight(m.windowHeight)
	return lipgloss.JoinHorizontal(lipgloss.Top, m.issues.View(), m.editor.View())
}

func formatJSONFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return err.Error()
	}
	var indented bytes.Buffer
	if err := json.Indent(&indented, data, "", "  "); err != nil {
		return string(data)
	}
	var highlighted bytes.Buffer
	if err := quick.Highlight(&highlighted, indented.String(), "json", "terminal256", "github"); err != nil {
		return indented.String()
	}
	return highlighted.String()
}
