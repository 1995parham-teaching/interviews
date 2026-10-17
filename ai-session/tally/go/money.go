// Package tally keeps a shared expense ledger for a group of people.
//
// Amounts are kept as integer cents so that every sum is exact.
package tally

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var amountPattern = regexp.MustCompile(`^\d+(\.\d+)?$`)

// ParseAmount parses a decimal amount such as "12.50" or "7" into cents.
func ParseAmount(text string) (int64, error) {
	if !amountPattern.MatchString(text) {
		return 0, fmt.Errorf("invalid amount: %s", text)
	}
	var cents int64
	if whole, fraction, ok := strings.Cut(text, "."); ok {
		w, _ := strconv.ParseInt(whole, 10, 64)
		f, _ := strconv.ParseInt(fraction, 10, 64)
		cents = w*100 + f
	} else {
		w, _ := strconv.ParseInt(text, 10, 64)
		cents = w * 100
	}
	if cents == 0 {
		return 0, fmt.Errorf("amount must be positive: %s", text)
	}
	return cents, nil
}

// FormatAmount formats cents as a decimal amount.
func FormatAmount(cents int64) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
}

// FormatSigned is FormatAmount with a "+" in front of positive values.
func FormatSigned(cents int64) string {
	if cents > 0 {
		return "+" + FormatAmount(cents)
	}
	return FormatAmount(cents)
}
