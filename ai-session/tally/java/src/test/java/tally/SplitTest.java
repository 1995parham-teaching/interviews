package tally;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.util.Map;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

class SplitTest {
    @Test
    void equal() {
        assertEquals(Map.of("ali", 3000L, "sara", 3000L, "reza", 3000L), Split.parse("equal:ali,sara,reza", 9000));
    }

    @Test
    void shares() {
        assertEquals(Map.of("ali", 1000L, "sara", 2000L), Split.parse("shares:ali=1,sara=2", 3000));
    }

    @Test
    void exact() {
        assertEquals(Map.of("ali", 2000L, "sara", 3000L), Split.parse("exact:ali=20.00,sara=30.00", 5000));
    }

    @Test
    void exactMoreThanAmount() {
        assertThrows(IllegalArgumentException.class, () -> Split.parse("exact:ali=20.00,sara=40.00", 5000));
    }

    @ParameterizedTest
    @ValueSource(strings = {"equal", "equal:", "halves:ali,sara", "shares:ali", "shares:ali=0"})
    void invalid(String spec) {
        assertThrows(IllegalArgumentException.class, () -> Split.parse(spec, 1000));
    }
}
