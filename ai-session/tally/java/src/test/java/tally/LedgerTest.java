package tally;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;

class LedgerTest {
    private static Ledger trip() {
        Ledger ledger = new Ledger();
        for (String name : List.of("ali", "sara", "reza")) {
            ledger.addMember(name);
        }
        return ledger;
    }

    private static void hotel(Ledger ledger) {
        ledger.addExpense("ali", 9000, Map.of("ali", 3000L, "sara", 3000L, "reza", 3000L), "hotel");
    }

    @Test
    void expenseUpdatesBalances() {
        Ledger ledger = trip();
        hotel(ledger);
        assertEquals(Map.of("ali", 6000L, "sara", -3000L, "reza", -3000L), ledger.balances());
    }

    @Test
    void paymentUpdatesBalances() {
        Ledger ledger = trip();
        hotel(ledger);
        ledger.recordPayment("sara", "ali", 3000);
        assertEquals(Map.of("ali", 3000L, "sara", 0L, "reza", -3000L), ledger.balances());
    }

    @Test
    void deleteRevertsExpense() {
        Ledger ledger = trip();
        hotel(ledger);
        Ledger.Expense dinner = ledger.addExpense("sara", 2000, Map.of("ali", 1000L, "sara", 1000L), "dinner");
        ledger.deleteExpense(dinner.id());
        assertEquals(Map.of("ali", 6000L, "sara", -3000L, "reza", -3000L), ledger.balances());
        assertEquals(List.of("hotel"), ledger.expenses().stream().map(Ledger.Expense::description).toList());
    }

    @Test
    void unknownMember() {
        Ledger ledger = trip();
        assertThrows(LedgerException.class, () -> ledger.addExpense("mina", 1000, Map.of("ali", 1000L), "coffee"));
        assertThrows(LedgerException.class, () -> ledger.recordPayment("ali", "mina", 1000));
    }
}
