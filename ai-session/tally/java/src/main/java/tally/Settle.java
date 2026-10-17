package tally;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/** Works out who should pay whom so that every balance goes back to zero. */
public final class Settle {
    /** One payment in a settlement plan. */
    public record Transfer(String sender, String receiver, long amount) {}

    private Settle() {}

    /**
     * Repeatedly lets the member who owes the most pay the member who is owed the most. Ties are
     * broken by name so that the plan is the same on every run.
     */
    public static List<Transfer> plan(Map<String, Long> balances) {
        List<String> names = balances.keySet().stream().sorted().toList();
        List<Transfer> transfers = new ArrayList<>();
        while (!names.isEmpty()) {
            String creditor = names.get(0);
            String debtor = names.get(0);
            for (String name : names) {
                if (balances.get(name) > balances.get(creditor)) {
                    creditor = name;
                }
                if (balances.get(name) < balances.get(debtor)) {
                    debtor = name;
                }
            }
            if (balances.get(creditor) <= 0 || balances.get(debtor) >= 0) {
                break;
            }
            long amount = Math.min(balances.get(creditor), -balances.get(debtor));
            balances.merge(creditor, -amount, Long::sum);
            balances.merge(debtor, amount, Long::sum);
            transfers.add(new Transfer(debtor, creditor, amount));
        }
        return transfers;
    }
}
