package tally

import (
	"maps"
	"slices"
	"testing"
)

func trip(t *testing.T) *Ledger {
	t.Helper()
	ledger := NewLedger()
	for _, name := range []string{"ali", "sara", "reza"} {
		if err := ledger.AddMember(name); err != nil {
			t.Fatal(err)
		}
	}
	return ledger
}

func hotel(t *testing.T, ledger *Ledger) {
	t.Helper()
	shares := map[string]int64{"ali": 3000, "sara": 3000, "reza": 3000}
	if _, err := ledger.AddExpense("ali", 9000, shares, "hotel"); err != nil {
		t.Fatal(err)
	}
}

func assertBalances(t *testing.T, ledger *Ledger, want map[string]int64) {
	t.Helper()
	if got := ledger.Balances(); !maps.Equal(got, want) {
		t.Errorf("balances = %v, want %v", got, want)
	}
}

func TestExpenseUpdatesBalances(t *testing.T) {
	ledger := trip(t)
	hotel(t, ledger)
	assertBalances(t, ledger, map[string]int64{"ali": 6000, "sara": -3000, "reza": -3000})
}

func TestPaymentUpdatesBalances(t *testing.T) {
	ledger := trip(t)
	hotel(t, ledger)
	if _, err := ledger.RecordPayment("sara", "ali", 3000); err != nil {
		t.Fatal(err)
	}
	assertBalances(t, ledger, map[string]int64{"ali": 3000, "sara": 0, "reza": -3000})
}

func TestDeleteRevertsExpense(t *testing.T) {
	ledger := trip(t)
	hotel(t, ledger)
	dinner, err := ledger.AddExpense("sara", 2000, map[string]int64{"ali": 1000, "sara": 1000}, "dinner")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.DeleteExpense(dinner.ID); err != nil {
		t.Fatal(err)
	}
	assertBalances(t, ledger, map[string]int64{"ali": 6000, "sara": -3000, "reza": -3000})
	var descriptions []string
	for _, e := range ledger.Expenses() {
		descriptions = append(descriptions, e.Description)
	}
	if !slices.Equal(descriptions, []string{"hotel"}) {
		t.Errorf("expenses = %v, want [hotel]", descriptions)
	}
}

func TestUnknownMember(t *testing.T) {
	ledger := trip(t)
	if _, err := ledger.AddExpense("mina", 1000, map[string]int64{"ali": 1000}, "coffee"); err == nil {
		t.Error("expense paid by an unknown member was accepted")
	}
	if _, err := ledger.RecordPayment("ali", "mina", 1000); err == nil {
		t.Error("payment to an unknown member was accepted")
	}
}

func TestPlan(t *testing.T) {
	got := Plan(map[string]int64{"ali": 3600, "sara": -1200, "reza": -2400})
	want := []Transfer{{"reza", "ali", 2400}, {"sara", "ali", 1200}}
	if !slices.Equal(got, want) {
		t.Errorf("Plan = %v, want %v", got, want)
	}
}

func TestNothingToSettle(t *testing.T) {
	if got := Plan(map[string]int64{"ali": 0, "sara": 0}); len(got) != 0 {
		t.Errorf("Plan = %v, want no transfers", got)
	}
}
