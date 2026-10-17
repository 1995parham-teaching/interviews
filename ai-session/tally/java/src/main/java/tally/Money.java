package tally;

import java.util.regex.Pattern;

/** Amounts are kept as integer cents so that every sum is exact. */
public final class Money {
    private static final Pattern AMOUNT = Pattern.compile("[0-9]+(\\.[0-9]+)?");

    private Money() {}

    /** Parses a decimal amount such as "12.50" or "7" into cents. */
    public static long parse(String text) {
        if (!AMOUNT.matcher(text).matches()) {
            throw new IllegalArgumentException("invalid amount: " + text);
        }
        long cents;
        int dot = text.indexOf('.');
        if (dot >= 0) {
            long whole = Long.parseLong(text.substring(0, dot));
            long fraction = Long.parseLong(text.substring(dot + 1));
            cents = whole * 100 + fraction;
        } else {
            cents = Long.parseLong(text) * 100;
        }
        if (cents == 0) {
            throw new IllegalArgumentException("amount must be positive: " + text);
        }
        return cents;
    }

    /** Formats cents as a decimal amount. */
    public static String format(long cents) {
        String sign = cents < 0 ? "-" : "";
        cents = Math.abs(cents);
        return String.format("%s%d.%02d", sign, cents / 100, cents % 100);
    }

    /** Like {@link #format}, with a "+" in front of positive values. */
    public static String formatSigned(long cents) {
        return cents > 0 ? "+" + format(cents) : format(cents);
    }
}
