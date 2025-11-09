package models

import (
	"time"

	"gorm.io/gorm"
)

// User and rewards
type User struct {
	gorm.Model
	Name string
}

// Rewards are source of truth for holdings
type Reward struct {
	gorm.Model
	UserID    uint
	SymbolID  uint
	Quantity  float64
	Timestamp time.Time
}

// Symbol and price history
type Symbol struct {
	gorm.Model
	Name string
}

type SymbolPriceHistory struct {
	gorm.Model
	SymbolID uint
	Price    float64
	TimeHour uint      // 0 to 23, ASSUME Market is open all day everyday
	Date     time.Time `gorm:"type:date"` // Date only
}

// Double entry ledger
type Account struct {
	gorm.Model
	Name        string
	Type        string
	Description string
}
type Transaction struct {
	gorm.Model
	Description string
	Entries     []Entry
}
type Entry struct {
	gorm.Model
	TransactionID uint
	AccountID     uint
	Type          string // "credit" or "debit"
	Amount        float64
}

func SeedDatabase(db *gorm.DB, clearData bool) {
	db.AutoMigrate(&User{}, &Reward{}, &Symbol{}, &SymbolPriceHistory{}, &Account{}, &Transaction{}, &Entry{})

	// Seed symbol, account once
	var symbolCount int64
	db.Model(Symbol{}).Count(&symbolCount)
	if symbolCount == 0 {
		db.Create([]Symbol{
			{Name: "RELIANCE"},
			{Name: "HDFCBANK"},
			{Name: "BSE"},
			{Name: "SBIN"},
			{Name: "ICICIBANK"},
			{Name: "TCS"},
			{Name: "BHARTIIARTL"},
			{Name: "LTF"},
			{Name: "BAJFINANCE"},
			{Name: "INFY"},
			{Name: "TMPV"},
			{Name: "TATASTEEL"},
			{Name: "AXISBANK"},
			{Name: "BEL"},
		})
	}
	var accountsCount int64
	db.Model(Account{}).Count(&accountsCount)
	if accountsCount == 0 {
		db.Create([]Account{
			{Name: "Cash", Type: "asset", Description: "This account represents company's liquid money"},
			{Name: "StockInvestments", Type: "asset", Description: "Represents the stocks the company Stocky currently owns"},
			{Name: "TransactionFees", Type: "expense", Description: "Brokerage, taxes, and all other fees incurred while trading"},
			{Name: "RewardExpense", Type: "expense", Description: "The value of free stocks given to users as part of promotions or loyalty programs"},
		})
	}

	if clearData {
		db.Where("1 = 1").Delete(&User{})
		db.Where("1 = 1").Delete(&Reward{})
		db.Where("1 = 1").Delete(&SymbolPriceHistory{})
		db.Where("1 = 1").Delete(&Transaction{})
		db.Where("1 = 1").Delete(&Entry{})
	}
}
