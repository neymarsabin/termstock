package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/neymarsabin/termstock/database"
	"github.com/neymarsabin/termstock/nepse"
	"gorm.io/gorm"
)

type model struct {
	symbols     []string
	quotes      map[string]nepse.Quote
	err         error
	loading     bool
	spinner     spinner.Model
	inputSymbol string
	addMode     bool
	textInput   textinput.Model
	db          *gorm.DB
}

func main() {
	p := tea.NewProgram(initProgram(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}

func initProgram() model {
	db := database.Open()
	symbolsData := database.SymbolsFromDb(db)

	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := model{
		symbols:   symbolsData,
		quotes:    make(map[string]nepse.Quote),
		loading:   true,
		spinner:   s,
		textInput: ti,
		db:        db,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchQuotes(m.symbols))
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if !m.addMode {
				return m, tea.Quit
			}

		case "a":
			if !m.addMode {
				m.addMode = true
				m.inputSymbol = ""
				return m, textinput.Blink
			}

		case "r":
			if !m.addMode {
				m.loading = true
				return m, fetchQuotes(m.symbols)
			}

		case "enter":
			// TODO: save the symbol in the database
			// TODO: set m.addMode = false
			if m.addMode {
				m.addMode = false
				m.inputSymbol = m.textInput.Value()

				if m.inputSymbol == "" {
					m.loading = true
					return m, fetchQuotes(m.symbols)
				}

				m.symbols = append(m.symbols, m.inputSymbol)
				m.loading = true
				_ = database.AddSymbol(m.inputSymbol, m.db)
				return m, fetchQuotes(m.symbols)
			}
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

	m.textInput, _ = m.textInput.Update(msg)

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
			Foreground(lipgloss.Color("#4287f5")).
			Render(m.spinner.View() + "Loading stock prices...")
	}

	if m.addMode {
		inputView := fmt.Sprintf(
			"\n\n\n Add the symbol to fetch quotes? \n\n%s\n\n",
			m.textInput.View(),
		) + "\n"

		return lipgloss.NewStyle().
			Align(lipgloss.Center).
			Height(40).
			Width(100).
			Render(menu + inputView)
	}

	var rows []string
	for _, symbol := range m.symbols {
		priceStyle := lipgloss.NewStyle()
		if m.quotes[symbol].Positive {
			priceStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#008000"))
		} else {
			priceStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
		}
		rows = append(rows, fmt.Sprintf("%s: %s", symbol, priceStyle.Render(fmt.Sprintf("%v (%v)", m.quotes[symbol].Price, m.quotes[symbol].PercentageChange))))
	}

	mainView := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Height(40).
		Width(100).
		Render(menu + "\n\n\n\n" + lipgloss.JoinVertical(lipgloss.Center, rows...))

	return mainView
}
