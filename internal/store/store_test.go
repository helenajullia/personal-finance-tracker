package store

import (
	"sync"
	"testing"
	"time"

	"atad-project/internal/money"
)

func day(s string) time.Time {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return d
}

func amount(s string) money.Amount {
	a, err := money.Parse(s)
	if err != nil {
		panic(err)
	}
	return a
}

func TestAddAssignsSequentialIDsFromOne(t *testing.T) {
	s := New()
	for want := 1; want <= 3; want++ {
		got := s.Add(day("2026-10-01"), amount("-1.00"), "item", "")
		if got.ID != want {
			t.Errorf("ID = %d, want %d", got.ID, want)
		}
	}
}

func TestAddKeepsFields(t *testing.T) {
	s := New()
	got := s.Add(day("2026-10-05"), amount("-12.50"), "groceries", "food")

	if !got.Date.Equal(day("2026-10-05")) {
		t.Errorf("Date = %v", got.Date)
	}
	if got.Amount.String() != "-12.50" {
		t.Errorf("Amount = %s, want -12.50", got.Amount)
	}
	if got.Description != "groceries" || got.Category != "food" {
		t.Errorf("Description/Category = %q/%q", got.Description, got.Category)
	}
}

func TestListOrdering(t *testing.T) {
	tests := []struct {
		name  string
		dates []string
		want  []int
	}{
		{"newest first", []string{"2026-10-01", "2026-10-03", "2026-10-02"}, []int{2, 3, 1}},
		{"same date by ascending id", []string{"2026-10-01", "2026-10-01", "2026-10-01"}, []int{1, 2, 3}},
		{"mixed ties", []string{"2026-10-01", "2026-10-02", "2026-10-02", "2026-10-01"}, []int{2, 3, 1, 4}},
		{"single", []string{"2026-10-01"}, []int{1}},
		{"empty", nil, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			for _, d := range tt.dates {
				s.Add(day(d), amount("1.00"), "x", "")
			}
			got := s.List()
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("position %d: ID = %d, want %d", i, got[i].ID, id)
				}
			}
		})
	}
}

func TestListReturnsCopy(t *testing.T) {
	s := New()
	s.Add(day("2026-10-01"), amount("1.00"), "original", "")

	list := s.List()
	list[0].Description = "changed"

	if s.List()[0].Description != "original" {
		t.Error("modifying the returned slice changed the store")
	}
}

func TestResetClearsAndRestartsIDs(t *testing.T) {
	s := New()
	s.Add(day("2026-10-01"), amount("1.00"), "a", "")
	s.Add(day("2026-10-02"), amount("2.00"), "b", "")

	s.Reset()

	if n := len(s.List()); n != 0 {
		t.Fatalf("after Reset, len = %d, want 0", n)
	}
	if id := s.Add(day("2026-10-03"), amount("3.00"), "c", "").ID; id != 1 {
		t.Errorf("first ID after Reset = %d, want 1", id)
	}
}

func TestConcurrentAdd(t *testing.T) {
	s := New()
	const n = 200

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Add(day("2026-10-01"), amount("1.00"), "x", "")
		}()
	}
	wg.Wait()

	list := s.List()
	if len(list) != n {
		t.Fatalf("len = %d, want %d", len(list), n)
	}
	seen := make(map[int]bool, n)
	for _, tr := range list {
		if tr.ID < 1 || tr.ID > n || seen[tr.ID] {
			t.Fatalf("bad or duplicate ID %d", tr.ID)
		}
		seen[tr.ID] = true
	}
}
