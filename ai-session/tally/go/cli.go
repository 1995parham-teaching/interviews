package tally

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

const usage = `commands:
  member <name>
  expense <payer> <amount> <split> <description>
  pay <sender> <receiver> <amount>
  delete <expense-id>
  balances
  settle`

// Run reads commands, one per line, runs them against a single ledger and
// returns the number of commands that failed.
func Run(in io.Reader, out io.Writer) (int, error) {
	ledger := NewLedger()
	failures := 0
	scanner := bufio.NewScanner(in)
	for number := 1; scanner.Scan(); number++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := Execute(ledger, line, out); err != nil {
			fmt.Fprintf(out, "error: line %d: %v\n", number, err)
			failures++
		}
	}
	return failures, scanner.Err()
}

func Execute(ledger *Ledger, line string, out io.Writer) error {
	fields := strings.Fields(line)
	command, args := fields[0], fields[1:]
	switch command {
	case "member":
		if err := arity(args, 1, "member <name>"); err != nil {
			return err
		}
		return ledger.AddMember(args[0])
	case "expense":
		if len(args) < 4 {
			return errors.New("usage: expense <payer> <amount> <split> <description>")
		}
		amount, err := ParseAmount(args[1])
		if err != nil {
			return err
		}
		shares, err := ParseSplit(args[2], amount)
		if err != nil {
			return err
		}
		expense, err := ledger.AddExpense(args[0], amount, shares, strings.Join(args[3:], " "))
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "#%d %s: %s paid %s\n", expense.ID, expense.Description, expense.Payer, FormatAmount(expense.Amount))
	case "pay":
		if err := arity(args, 3, "pay <sender> <receiver> <amount>"); err != nil {
			return err
		}
		amount, err := ParseAmount(args[2])
		if err != nil {
			return err
		}
		payment, err := ledger.RecordPayment(args[0], args[1], amount)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s paid %s %s\n", payment.Sender, payment.Receiver, FormatAmount(payment.Amount))
	case "delete":
		if err := arity(args, 1, "delete <expense-id>"); err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil || id < 0 || strings.HasPrefix(args[0], "+") {
			return fmt.Errorf("invalid expense id: %s", args[0])
		}
		expense, err := ledger.DeleteExpense(id)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "deleted #%d %s\n", expense.ID, expense.Description)
	case "balances":
		if err := arity(args, 0, "balances"); err != nil {
			return err
		}
		printBalances(ledger.Balances(), out)
	case "settle":
		if err := arity(args, 0, "settle"); err != nil {
			return err
		}
		transfers := Plan(ledger.Balances())
		if len(transfers) == 0 {
			fmt.Fprintln(out, "nothing to settle")
		}
		for _, t := range transfers {
			fmt.Fprintf(out, "%s -> %s %s\n", t.Sender, t.Receiver, FormatAmount(t.Amount))
		}
	default:
		return fmt.Errorf("unknown command: %s\n%s", command, usage)
	}
	return nil
}

func printBalances(balances map[string]int64, out io.Writer) {
	names := make([]string, 0, len(balances))
	width := 0
	for name := range balances {
		names = append(names, name)
		width = max(width, len(name))
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(out, "%-*s  %s\n", width, name, FormatSigned(balances[name]))
	}
}

func arity(args []string, n int, usage string) error {
	if len(args) != n {
		return fmt.Errorf("usage: %s", usage)
	}
	return nil
}
