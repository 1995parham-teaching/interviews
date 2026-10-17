package tally;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

class MoneyTest {
    @Test
    void parsesWholeAndCents() {
        assertEquals(1250, Money.parse("12.50"));
        assertEquals(700, Money.parse("7"));
        assertEquals(5, Money.parse("0.05"));
    }

    @ParameterizedTest
    @ValueSource(strings = {"", "abc", "-3", "1,5", "0"})
    void rejectsGarbage(String text) {
        assertThrows(IllegalArgumentException.class, () -> Money.parse(text));
    }

    @Test
    void formats() {
        assertEquals("12.50", Money.format(1250));
        assertEquals("0.05", Money.format(5));
        assertEquals("-3.00", Money.format(-300));
    }

    @Test
    void formatsSigned() {
        assertEquals("+40.00", Money.formatSigned(4000));
        assertEquals("0.00", Money.formatSigned(0));
    }
}
