package tally;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Holds a group's members, expenses and payments.
 *
 * <p>Balances are kept up to date as entries are recorded. A positive balance means the group owes
 * that member money; a negative one means the member owes the group.
 */
public final class Ledger {
    /** Money one member paid on behalf of others. */
    public record Expense(int id, String payer, long amount, Map<String, Long> shares, String description) {}

    /** Money one member gave directly to another. */
    public record Payment(String sender, String receiver, long amount) {}

    private final List<String> members = new ArrayList<>();
    private final List<Expense> expenses = new ArrayList<>();
    private final List<Payment> payments = new ArrayList<>();
    private final Map<String, Long> balances = new LinkedHashMap<>();

    public void addMember(String name) {
        if (balances.containsKey(name)) {
            throw new LedgerException("member already exists: " + name);
        }
        members.add(name);
        balances.put(name, 0L);
    }

    public List<String> members() {
        return List.copyOf(members);
    }

    public Expense addExpense(String payer, long amount, Map<String, Long> shares, String description) {
        require(payer);
        shares.keySet().forEach(this::require);
        Expense expense = new Expense(expenses.size() + 1, payer, amount, shares, description);
        expenses.add(expense);
        apply(expense, 1);
        return expense;
    }

    public Expense deleteExpense(int id) {
        for (Expense expense : expenses) {
            if (expense.id() == id) {
                expenses.remove(expense);
                apply(expense, -1);
                return expense;
            }
        }
        throw new LedgerException("no expense #" + id);
    }

    public List<Expense> expenses() {
        return List.copyOf(expenses);
    }

    public Payment recordPayment(String sender, String receiver, long amount) {
        require(sender);
        require(receiver);
        if (sender.equals(receiver)) {
            throw new LedgerException(sender + " cannot pay themselves");
        }
        Payment payment = new Payment(sender, receiver, amount);
        payments.add(payment);
        balances.merge(sender, amount, Long::sum);
        balances.merge(receiver, -amount, Long::sum);
        return payment;
    }

    public List<Payment> payments() {
        return List.copyOf(payments);
    }

    public Map<String, Long> balances() {
        return balances;
    }

    private void apply(Expense expense, long sign) {
        balances.merge(expense.payer(), sign * expense.amount(), Long::sum);
        expense.shares().forEach((name, share) -> balances.merge(name, -sign * share, Long::sum));
    }

    private void require(String name) {
        if (!balances.containsKey(name)) {
            throw new LedgerException("unknown member: " + name);
        }
    }
}
