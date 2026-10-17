package tally

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseSplit turns a split specification into the amount each member owes for
// an expense. A specification has a kind and a list of members:
//
//	equal:ali,sara,reza            everyone owes the same
//	shares:ali=1,sara=2            owed in proportion to the weights
//	exact:ali=20.00,sara=30.00     the amounts are given directly
func ParseSplit(spec string, amount int64) (map[string]int64, error) {
	kind, rest, ok := strings.Cut(spec, ":")
	if !ok || rest == "" {
		return nil, fmt.Errorf("invalid split: %s", spec)
	}
	switch kind {
	case "equal":
		return Equal(amount, strings.Split(rest, ",")), nil
	case "shares":
		weights, err := parsePairs(rest, parseWeight)
		if err != nil {
			return nil, err
		}
		return ByShares(amount, weights), nil
	case "exact":
		parts, err := parsePairs(rest, ParseAmount)
		if err != nil {
			return nil, err
		}
		return Exact(amount, parts)
	default:
		return nil, fmt.Errorf("unknown split kind: %s", kind)
	}
}

func Equal(amount int64, names []string) map[string]int64 {
	each := amount / int64(len(names))
	shares := make(map[string]int64, len(names))
	for _, name := range names {
		shares[name] = each
	}
	return shares
}

func ByShares(amount int64, weights map[string]int64) map[string]int64 {
	var total int64
	for _, w := range weights {
		total += w
	}
	shares := make(map[string]int64, len(weights))
	for name, w := range weights {
		shares[name] = amount * w / total
	}
	return shares
}

func Exact(amount int64, parts map[string]int64) (map[string]int64, error) {
	var total int64
	for _, a := range parts {
		total += a
	}
	if total > amount {
		return nil, fmt.Errorf("exact parts add up to %s, more than %s", FormatAmount(total), FormatAmount(amount))
	}
	return parts, nil
}

func parsePairs(text string, parse func(string) (int64, error)) (map[string]int64, error) {
	pairs := make(map[string]int64)
	for _, item := range strings.Split(text, ",") {
		name, value, ok := strings.Cut(item, "=")
		if !ok || name == "" || value == "" {
			return nil, fmt.Errorf("expected name=value, got: %s", item)
		}
		v, err := parse(value)
		if err != nil {
			return nil, err
		}
		pairs[name] = v
	}
	return pairs, nil
}

func parseWeight(text string) (int64, error) {
	w, err := strconv.ParseInt(text, 10, 64)
	if err != nil || w <= 0 || strings.HasPrefix(text, "+") {
		return 0, fmt.Errorf("invalid weight: %s", text)
	}
	return w, nil
}
