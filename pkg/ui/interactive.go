// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package ui

import (
	"fmt"
	"strings"

	"github.com/chainguard-dev/bomshell/pkg/render"
	"github.com/chainguard-dev/bomshell/pkg/shell"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/cel-go/cel"
)

const Prompt = "🐚❯ "

var titleStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
}()

/*
	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
*/type historyEntry struct {
	expression string
	result     string
	isError    bool
}

type History []historyEntry

func (h *History) Append(entry historyEntry) {
	*h = append(*h, entry)
}

type model struct {
	ready    bool
	bomshell *shell.Bomshell
	debug    bool
	renderer render.Renderer
	history  History
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
		case tea.KeyEnter: // tea.KeyUp
			// Execute the expression:
			result, err := m.bomshell.Run(m.textarea.Value())
			if err == nil {
				m.history = append(m.history, historyEntry{
					m.textarea.Value(),
					m.renderer.Display(result),
					false,
				})
			} else {
				m.history = append(m.history, historyEntry{
					m.textarea.Value(),
					err.Error(),
					true,
				})
			}

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

func initModel(bomshell *shell.Bomshell) model {
	ta := textarea.New()
	ta.Placeholder = "type a bomshell expression..."
	ta.Focus()

	ta.Prompt = Prompt
	ta.SetHeight(1)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		bomshell: bomshell,
		history:  []historyEntry{},
		renderer: render.NewTTY(),
		textarea: ta,
	}
}

type Interactive struct {
	ui       *tea.Program
	bomshell *shell.Bomshell
}

func (i *Interactive) Start() error {
	if _, err := i.ui.Run(); err != nil {
		return fmt.Errorf("starting UI: %w", err)
	}
	return nil
}

func NewInteractive(opts shell.Options) (*Interactive, error) {
	opts.EnvOptions = append(opts.EnvOptions, cel.Lib(InteractiveSubshell{}))
	bomshell, err := shell.NewWithOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("creating bomshell environment: %w", err)
	}
	/*
		if err := bomshell.RunFile(program); err != nil {
			logrus.Fatal(err)
		}
	*/
	return &Interactive{
		bomshell: bomshell,
		ui: tea.NewProgram(
			initModel(bomshell),
			tea.WithAltScreen(),
			// tea.WithMouseCellMotion(),
		),
	}, nil
}
