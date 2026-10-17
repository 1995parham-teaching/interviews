package tally

import "fmt"

// Expense is money one member paid on behalf of others.
type Expense struct {
	ID          int
	Payer       string
	Amount      int64
	Shares      map[string]int64
	Description string
}

// Payment is money one member gave directly to another.
type Payment struct {
	Sender   string
	Receiver string
	Amount   int64
}

// Ledger holds a group's members, expenses and payments.
//
// Balances are kept up to date as entries are recorded. A positive balance
// means the group owes that member money; a negative one means the member owes
// the group.
type Ledger struct {
	members  []string
	expenses []Expense
	payments []Payment
	balances map[string]int64
}

func NewLedger() *Ledger {
	return &Ledger{balances: make(map[string]int64)}
}

func (l *Ledger) AddMember(name string) error {
	if _, ok := l.balances[name]; ok {
		return fmt.Errorf("member already exists: %s", name)
	}
	l.members = append(l.members, name)
	l.balances[name] = 0
	return nil
}

func (l *Ledger) Members() []string {
	return append([]string(nil), l.members...)
}

func (l *Ledger) AddExpense(payer string, amount int64, shares map[string]int64, description string) (Expense, error) {
	names := []string{payer}
	for name := range shares {
		names = append(names, name)
	}
	if err := l.require(names...); err != nil {
		return Expense{}, err
	}
	expense := Expense{
		ID:          len(l.expenses) + 1,
		Payer:       payer,
		Amount:      amount,
		Shares:      shares,
		Description: description,
	}
	l.expenses = append(l.expenses, expense)
	l.apply(expense, 1)
	return expense, nil
}

func (l *Ledger) DeleteExpense(id int) (Expense, error) {
	for i, expense := range l.expenses {
		if expense.ID == id {
			l.expenses = append(l.expenses[:i], l.expenses[i+1:]...)
			l.apply(expense, -1)
			return expense, nil
		}
	}
	return Expense{}, fmt.Errorf("no expense #%d", id)
}

func (l *Ledger) Expenses() []Expense {
	return append([]Expense(nil), l.expenses...)
}

func (l *Ledger) RecordPayment(sender, receiver string, amount int64) (Payment, error) {
	if err := l.require(sender, receiver); err != nil {
		return Payment{}, err
	}
	if sender == receiver {
		return Payment{}, fmt.Errorf("%s cannot pay themselves", sender)
	}
	payment := Payment{Sender: sender, Receiver: receiver, Amount: amount}
	l.payments = append(l.payments, payment)
	l.balances[sender] += amount
	l.balances[receiver] -= amount
	return payment, nil
}

func (l *Ledger) Payments() []Payment {
	return append([]Payment(nil), l.payments...)
}

func (l *Ledger) Balances() map[string]int64 {
	return l.balances
}

func (l *Ledger) apply(expense Expense, sign int64) {
	l.balances[expense.Payer] += sign * expense.Amount
	for name, share := range expense.Shares {
		l.balances[name] -= sign * share
	}
}

func (l *Ledger) require(names ...string) error {
	for _, name := range names {
		if _, ok := l.balances[name]; !ok {
			return fmt.Errorf("unknown member: %s", name)
		}
	}
	return nil
}
