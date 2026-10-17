package tally

import "sort"

// Transfer is one payment in a settlement plan.
type Transfer struct {
	Sender   string
	Receiver string
	Amount   int64
}

// Plan works out who should pay whom so that every balance goes back to zero.
// It repeatedly lets the member who owes the most pay the member who is owed
// the most. Ties are broken by name so that the plan is the same on every run.
func Plan(balances map[string]int64) []Transfer {
	names := make([]string, 0, len(balances))
	for name := range balances {
		names = append(names, name)
	}
	sort.Strings(names)

	var transfers []Transfer
	for len(names) > 0 {
		creditor, debtor := names[0], names[0]
		for _, name := range names {
			if balances[name] > balances[creditor] {
				creditor = name
			}
			if balances[name] < balances[debtor] {
				debtor = name
			}
		}
		if balances[creditor] <= 0 || balances[debtor] >= 0 {
			break
		}
		amount := min(balances[creditor], -balances[debtor])
		balances[creditor] -= amount
		balances[debtor] += amount
		transfers = append(transfers, Transfer{Sender: debtor, Receiver: creditor, Amount: amount})
	}
	return transfers
}
