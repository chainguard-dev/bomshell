package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type historyEntry struct {
	expression string
	result     string
}

type model struct {
	ready    bool
	history  []historyEntry
	viewport viewport.Model
	textarea textarea.Model
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		textareaCmd tea.Cmd
		cmd         tea.Cmd
		cmds        []tea.Cmd
	)

	if m.ready {
		m.textarea, textareaCmd = m.textarea.Update(msg)
		cmds = append(cmds, textareaCmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.history = append(m.history, historyEntry{
				"Val: " + m.textarea.Value(), "resultado",
			})

			// content := fmt.Sprintf("Va pues (len es %d):\n", len(m.history))
			content := ""
			for _, e := range m.history {
				content = content + "> " + e.expression + "\n"
				content = content + "- " + e.result + "\n\n"
			}
			m.viewport.SetContent(content)
			m.textarea.Reset()
			m.viewport.GotoBottom()
			cmds = append(cmds, viewport.Sync(m.viewport))
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		m.textarea.SetWidth(msg.Width)

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = true
			m.viewport.KeyMap = viewport.KeyMap{}
			m.ready = true
			m.textarea.Focus()
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		cmds = append(cmds, viewport.Sync(m.viewport))
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render("💣🐚 BOMSHELL v0.0.1")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	return m.textarea.View()
	/*
		info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
		line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
		return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
	*/
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func initModel() model {
	ta := textarea.New()
	ta.Placeholder = "type a bomshell expression..."
	ta.Focus()

	ta.Prompt = "🐚❯ "
	ta.SetHeight(1)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		history:  []historyEntry{},
		textarea: ta,
	}
}

func main() {
	p := tea.NewProgram(
		initModel(),
		tea.WithAltScreen(),
		//tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
