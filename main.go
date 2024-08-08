package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/neymarsabin/termstock/nepse"
)

type model struct {
	symbols     []string
	quotes      map[string]nepse.Quote
	err         error
	loading     bool
	spinner     spinner.Model
	inputSymbol string
	addingMode  bool
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return model{
		symbols: []string{"NABIL"},
		quotes:  make(map[string]nepse.Quote),
		loading: true,
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchQuotes())
}

type tickMsg struct{}
type quotesMsg map[string]nepse.Quote
type errMsg error

func fetchQuotes() tea.Cmd {
	return func() tea.Msg {
		nabilData := nepse.Scrape()
		quotes := map[string]nepse.Quote{
			"NABIL": *nabilData,
		}
		nepse.Scrape()
		return quotesMsg(quotes)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "a":

		case "r":
			m.loading = true
			return m, fetchQuotes()

		case "enter":
		}

	case tickMsg:
		return m, tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	case quotesMsg:
		m.quotes = msg
		m.loading = false
	case errMsg:
		m.err = msg
		m.loading = false
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.loading {
		cmds = append(cmds, m.spinner.Tick)
	} else {
		cmds = append(cmds, tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
			return tickMsg{}
		}))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	menuStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Foreground(lipgloss.Color("230"))

	menu := lipgloss.JoinHorizontal(
		lipgloss.Top,
		menuStyle.Render("[a] Add Symbol"),
		menuStyle.Render("[r] Refresh"),
		menuStyle.Render("[q] Quit"),
	)

	if m.loading {
		return lipgloss.NewStyle().
			Align(lipgloss.Center).
			Height(40).
			Width(100).
			Render(m.spinner.View() + "Loading stock prices...")
	}

	var rows []string
	for _, symbol := range m.symbols {
		priceStyle := lipgloss.NewStyle()
		if m.quotes[symbol].Positive {
			priceStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#008000"))
		} else {
			priceStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
		}
		rows = append(rows, fmt.Sprintf("%s: %s", symbol, priceStyle.Render(fmt.Sprintf("%v", m.quotes[symbol].Price))))
	}

	mainView := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Height(40).
		Width(100).
		Render(menu + "\n\n\n\n" + lipgloss.JoinVertical(lipgloss.Center, rows...))

	return mainView
}
