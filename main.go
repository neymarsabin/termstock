package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	symbols     []string
	quotes      map[string]float64
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
		symbols: []string{"AAPL", "GOOGL", "MSFT"},
		quotes:  make(map[string]float64),
		loading: true,
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchQuotes())
}

type tickMsg struct{}
type quotesMsg map[string]float64
type errMsg error

func fetchQuotes() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(3 * time.Second)
		quotes := map[string]float64{
			"AAPL":  145.09,
			"GOOGL": 2729.30,
			"MSFT":  289.67,
		}
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

	if m.loading {
		return lipgloss.NewStyle().
			Align(lipgloss.Center).
			Height(50).
			Width(50).
			Render(m.spinner.View() + " Loading stock prices...")
	}

	var rows []string
	for _, symbol := range m.symbols {
		priceStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
		rows = append(rows, fmt.Sprintf("%s: %s", symbol, priceStyle.Render(fmt.Sprintf("%.2f", m.quotes[symbol]))))
	}

	s := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Height(50).
		Width(50).
		Render("Stock Prices:\n\n" + lipgloss.JoinVertical(lipgloss.Center, rows...))

	return s + "\n\n Press q to quit."
}
