package store

import (
	"sort"
	"sync"
	"time"

	"atad-project/internal/money"
)

type Transaction struct {
	ID          int
	Date        time.Time
	Amount      money.Amount
	Description string
	Category    string
}

type Store struct {
	mu           sync.Mutex
	nextID       int
	transactions []Transaction
}

func New() *Store {
	return &Store{nextID: 1}
}

func (s *Store) Add(date time.Time, amount money.Amount, description, category string) Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := Transaction{
		ID:          s.nextID,
		Date:        date,
		Amount:      amount,
		Description: description,
		Category:    category,
	}
	s.nextID++
	s.transactions = append(s.transactions, t)
	return t
}

func (s *Store) List() []Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Transaction, len(s.transactions))
	copy(out, s.transactions)
	sort.Slice(out, func(i, j int) bool {
		if !out[i].Date.Equal(out[j].Date) {
			return out[i].Date.After(out[j].Date)
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID = 1
	s.transactions = nil
}
