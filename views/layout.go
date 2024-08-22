package views

import (
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	paddingTop    = 1
	paddingRight  = 2
	paddingBottom = 1
	paddingLeft   = 2
)

type State struct {
	LoadingPage   bool
	AboutPage     bool
	DashboardPage bool
	BoxPage       bool
	AddPage       bool
}

type InputMode struct {
	Symbol    string
	TextInput textinput.Model
}

type Layout struct {
	height        int
	width         int
	paddingTop    int
	paddingRight  int
	paddingBottom int
	paddingLeft   int
}

var MenuItems = []string{
	"[a] Add Symbol",
	"[r] Refresh",
	"[q] Quit",
}

// Style definitions.
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render

	descStyle = lipgloss.NewStyle().MarginTop(1)

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)

	// menuStyle
	menuStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Foreground(lipgloss.Color("230"))

	// Page.
	pageStyle = lipgloss.NewStyle().Padding(paddingTop, paddingRight, paddingBottom, paddingLeft)
)

func PhysicalWidth() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal("Could not get terminal width and height: ")
	}

	return width, height
}

func MainHeader() string {
	headerString := strings.Builder{}
	var (
		title strings.Builder
	)

	desc := lipgloss.JoinVertical(lipgloss.Left,
		descStyle.Render("NEPSE stocks data in the terminal"),
		infoStyle.Render("by neymarsabin"+divider+url("https://github.com/neymarsabin/termstock")),
	)

	row := lipgloss.JoinHorizontal(lipgloss.Top, title.String(), desc)
	headerString.WriteString(row)

	return pageStyle.Render(headerString.String())
}

func LayoutView(children []string) string {
	docStyle := ""

	for _, c := range children {
		docStyle = docStyle + "\n" + c
	}

	return lipgloss.NewStyle().
		Height(60).
		Width(400).
		Render(MainHeader() + docStyle)
}

func SpinnerView(spinnerView string) string {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Height(50).
		Width(50).
		Foreground(lipgloss.Color("#4287f5")).
		Render(spinnerView + "Loading Stock Prices.....")
}

func MenuView() string {
	docStyle := []string{}
	for _, m := range MenuItems {
		docStyle = append(docStyle, menuStyle.Render(m))
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Right,
		docStyle...,
	)
}
