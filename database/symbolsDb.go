package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Symbols struct {
	ID     uint   `gorm:"primaryKey"`
	Symbol string `gorm:"type:varchar(100)"`
}

const dbFile = "termstock.db"

func Open() *gorm.DB {
	// Create database directory if it doesn't exist: $HOME/.config/batterarch/
	homeDir := getHomeDirectory()
	databasePath := fmt.Sprintf("%s/.config/termstock/", homeDir)

	if _, err := os.Stat(databasePath); os.IsNotExist(err) {
		os.MkdirAll(databasePath, 0700)
	}

	fullDatabasePath := fmt.Sprintf("%s%s", databasePath, dbFile)
	db, err := gorm.Open(sqlite.Open(fullDatabasePath), &gorm.Config{})
	if err != nil {
		fmt.Println("Error while initializing database: ", err)
		os.Exit(1)
	}

	db.AutoMigrate(&Symbols{})
	return db
}

func getHomeDirectory() string {
	home, _ := os.UserHomeDir()
	return home
}

func SymbolsFromDb(db *gorm.DB) []string {
	var allSymbols []Symbols
	var returnSymbols []string
	db.Find(&allSymbols)

	for _, r := range allSymbols {
		returnSymbols = append(returnSymbols, r.Symbol)
	}

	return returnSymbols
}

func AddSymbol(symbol string, db *gorm.DB) error {
	_ = db.Exec(`INSERT INTO symbols (symbol) VALUES (?)`, symbol)

	return nil
}
