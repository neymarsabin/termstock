package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/neymarsabin/termstock/nepse"
	"github.com/neymarsabin/termstock/views"
	"gorm.io/gorm"
)

type Symbols []string
type Quotes map[string]nepse.Quote

type Model struct {
	Symbols   Symbols
	Quotes    Quotes
	Err       error
	Db        *gorm.DB
	State     views.State
	InputMode views.InputMode
	Spinner   spinner.Model
}

func DefaultModel() Model {
	return Model{
		State: views.State{
			LoadingPage:   true,
			AboutPage:     false,
			DashboardPage: false,
			BoxPage:       false,
			AddPage:       false,
		},
		Symbols: []string{},
	}
}
