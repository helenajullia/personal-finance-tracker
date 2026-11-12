package money

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type Amount int64

var (
	ErrEmpty      = errors.New("amount is empty")
	ErrFormat     = errors.New("amount must be a decimal with exactly two places, like -12.50")
	ErrOutOfRange = errors.New("amount is out of range")
)

func Parse(s string) (Amount, error) {
	if s == "" {
		return 0, ErrEmpty
	}

	negative := false
	switch s[0] {
	case '-':
		negative = true
		s = s[1:]
	case '+':
		s = s[1:]
	}

	whole, frac, found := strings.Cut(s, ".")
	if !found || len(frac) != 2 || whole == "" {
		return 0, ErrFormat
	}
	if !allDigits(whole) || !allDigits(frac) {
		return 0, ErrFormat
	}

	var cents int64
	for _, r := range whole {
		d := int64(r - '0')
		if cents > (math.MaxInt64-d)/10 {
			return 0, ErrOutOfRange
		}
		cents = cents*10 + d
	}
	if cents > (math.MaxInt64-int64(frac[0]-'0')*10-int64(frac[1]-'0'))/100 {
		return 0, ErrOutOfRange
	}
	cents = cents*100 + int64(frac[0]-'0')*10 + int64(frac[1]-'0')

	if negative {
		cents = -cents
	}
	return Amount(cents), nil
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func (a Amount) String() string {
	sign := ""
	v := int64(a)
	if v < 0 {
		sign = "-"
		v = -v
	}
	return fmt.Sprintf("%s%d.%02d", sign, v/100, v%100)
}

func (a Amount) Add(b Amount) (Amount, error) {
	sum := int64(a) + int64(b)
	if (int64(b) > 0 && sum < int64(a)) || (int64(b) < 0 && sum > int64(a)) {
		return 0, ErrOutOfRange
	}
	return Amount(sum), nil
}

func (a Amount) Cents() int64 {
	return int64(a)
}
