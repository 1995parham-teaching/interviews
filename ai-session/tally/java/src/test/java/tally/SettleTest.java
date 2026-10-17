package tally;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;

class SettleTest {
    @Test
    void plan() {
        Map<String, Long> balances = new HashMap<>(Map.of("ali", 3600L, "sara", -1200L, "reza", -2400L));
        assertEquals(
                List.of(new Settle.Transfer("reza", "ali", 2400), new Settle.Transfer("sara", "ali", 1200)),
                Settle.plan(balances));
    }

    @Test
    void nothingToSettle() {
        assertEquals(List.of(), Settle.plan(new HashMap<>(Map.of("ali", 0L, "sara", 0L))));
    }
}
