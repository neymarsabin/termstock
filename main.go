package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/neymarsabin/termstock/database"
	"github.com/neymarsabin/termstock/nepse"
	"github.com/neymarsabin/termstock/views"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	p := tea.NewProgram(initProgram(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}

func initProgram() Model {
	db := database.Open()
	symbolsData := database.SymbolsFromDb(db)
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "Enter a symbol"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	m := DefaultModel()
	m.Symbols = symbolsData
	m.Quotes = make(map[string]nepse.Quote)
	m.Db = db
	m.Spinner = s
	m.InputMode.TextInput = ti

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, fetchQuotes(m.Symbols))
}

type tickMsg struct{}
type quotesMsg map[string]nepse.Quote
type errMsg error

func fetchQuotes(symbols []string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(3 * time.Second)
		quotes := make(map[string]nepse.Quote)

		for _, symbol := range symbols {
			quoteData := nepse.ScrapeBySymbol(symbol)
			quotes[symbol] = *quoteData
		}

		return quotesMsg(quotes)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if !m.State.AddPage {
				return m, tea.Quit
			}

		case "a":
			if !m.State.AddPage {
				m.State.AddPage = true
				m.InputMode.Symbol = ""
				return m, textinput.Blink
			}

		case "r":
			if !m.State.AddPage {
				m.State.LoadingPage = true
				return m, fetchQuotes(m.Symbols)
			}

		case "enter":
			if m.State.AddPage {
				m.State.AddPage = false
				m.InputMode.Symbol = m.InputMode.TextInput.Value()

				if m.InputMode.Symbol == "" {
					m.State.LoadingPage = true
					return m, fetchQuotes(m.Symbols)
				}

				m.Symbols = append(m.Symbols, m.InputMode.Symbol)
				m.State.LoadingPage = true
				_ = database.AddSymbol(m.InputMode.Symbol, m.Db)
				return m, fetchQuotes(m.Symbols)
			}
		}

	case tickMsg:
		return m, tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
			return tickMsg{}
		})

	case quotesMsg:
		m.Quotes = Quotes(msg)
		m.State.LoadingPage = false

	case errMsg:
		m.Err = msg
		m.State.LoadingPage = false

	default:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.State.LoadingPage {
		cmds = append(cmds, m.Spinner.Tick)
	} else {
		cmds = append(cmds, tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
			return tickMsg{}
		}))
	}

	m.InputMode.TextInput, _ = m.InputMode.TextInput.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error: %v", m.Err)
	}

	menu := views.MenuView()

	if m.State.LoadingPage {
		return views.LayoutView([]string{views.SpinnerView(m.Spinner.View())})
	}

	if m.State.AddPage {
		inputView := fmt.Sprintf(
			"\n\n\n Add the symbol to fetch quotes? \n\n%s\n\n",
			m.InputMode.TextInput.View(),
		) + "\n"

		return lipgloss.NewStyle().
			Align(lipgloss.Center).
			Height(40).
			Width(100).
			Render(menu + inputView)
	}

	var rows []string
	for _, symbol := range m.Symbols {
		priceStyle := lipgloss.NewStyle()
		if m.Quotes[symbol].Positive {
			priceStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#008000"))
		} else {
			priceStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
		}
		rows = append(rows, fmt.Sprintf("%s: %s", symbol, priceStyle.Render(fmt.Sprintf("%v (%v)", m.Quotes[symbol].Price, m.Quotes[symbol].PercentageChange))))
	}

	symbolsList := lipgloss.NewStyle().
		PaddingTop(1).
		PaddingRight(2).
		PaddingBottom(1).
		PaddingLeft(2).
		Render(lipgloss.JoinVertical(lipgloss.Center, rows...))

	return views.LayoutView([]string{menu, symbolsList})
}
