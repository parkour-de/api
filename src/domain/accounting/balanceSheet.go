package accounting

import (
	"pkv/api/src/domain"
	"time"
)

type BalanceSheet struct {
	domain.Entity
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Date          time.Time `json:"date"`
	BalanceChange float64   `json:"balance_change"`
	Notes         string    `json:"notes"`
}
