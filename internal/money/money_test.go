package money

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Amount
		wantErr error
	}{
		{"positive", "12.50", 1250, nil},
		{"negative", "-12.50", -1250, nil},
		{"explicit plus", "+3.05", 305, nil},
		{"zero", "0.00", 0, nil},
		{"negative zero", "-0.00", 0, nil},
		{"under one", "0.07", 7, nil},
		{"large", "9999999.99", 999999999, nil},
		{"empty", "", 0, ErrEmpty},
		{"one decimal place", "12.5", 0, ErrFormat},
		{"three decimal places", "12.505", 0, ErrFormat},
		{"no decimal part", "12", 0, ErrFormat},
		{"decimal comma", "12,50", 0, ErrFormat},
		{"missing whole part", ".50", 0, ErrFormat},
		{"letters", "xx", 0, ErrFormat},
		{"exponent", "1e3.00", 0, ErrFormat},
		{"two signs", "--1.00", 0, ErrFormat},
		{"spaces", " 1.00", 0, ErrFormat},
		{"only sign", "-", 0, ErrFormat},
		{"overflow", "92233720368547758.08", 0, ErrOutOfRange},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Parse(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseMaxInt64(t *testing.T) {
	got, err := Parse("92233720368547758.07")
	if err != nil {
		t.Fatalf("Parse(max) error = %v", err)
	}
	if got != Amount(9223372036854775807) {
		t.Errorf("Parse(max) = %d", got)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input Amount
		want  string
	}{
		{"positive", 1250, "12.50"},
		{"negative", -1250, "-12.50"},
		{"zero", 0, "0.00"},
		{"single cent", 1, "0.01"},
		{"negative single cent", -1, "-0.01"},
		{"whole", 500, "5.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Errorf("Amount(%d).String() = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseStringRoundTrip(t *testing.T) {
	for _, s := range []string{"0.00", "0.01", "-0.01", "12.50", "-12.50", "1000000.00"} {
		a, err := Parse(s)
		if err != nil {
			t.Fatalf("Parse(%q): %v", s, err)
		}
		if got := a.String(); got != s {
			t.Errorf("round trip of %q gave %q", s, got)
		}
	}
}

func TestAddIsExact(t *testing.T) {
	a, _ := Parse("0.10")
	b, _ := Parse("0.20")
	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("Add error = %v", err)
	}
	if got := sum.String(); got != "0.30" {
		t.Errorf("0.10 + 0.20 = %q, want %q", got, "0.30")
	}
}

func TestAddManySmallAmounts(t *testing.T) {
	var total Amount
	one, _ := Parse("0.10")
	for i := 0; i < 1000; i++ {
		var err error
		total, err = total.Add(one)
		if err != nil {
			t.Fatalf("Add error = %v", err)
		}
	}
	if got := total.String(); got != "100.00" {
		t.Errorf("1000 x 0.10 = %q, want %q", got, "100.00")
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a, b    Amount
		want    Amount
		wantErr error
	}{
		{"positive and positive", 100, 250, 350, nil},
		{"positive and negative", 1000, -250, 750, nil},
		{"cancels out", 500, -500, 0, nil},
		{"both negative", -100, -200, -300, nil},
		{"positive overflow", Amount(9223372036854775807), 1, 0, ErrOutOfRange},
		{"negative overflow", Amount(-9223372036854775807 - 1), -1, 0, ErrOutOfRange},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Add(tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Add error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Add = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCents(t *testing.T) {
	a, _ := Parse("-12.50")
	if a.Cents() != -1250 {
		t.Errorf("Cents() = %d, want -1250", a.Cents())
	}
}
