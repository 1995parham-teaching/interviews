package tally;

import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.function.ToLongFunction;
import java.util.regex.Pattern;

/**
 * Turns a split specification into the amount each member owes for an expense. A specification
 * has a kind and a list of members:
 *
 * <pre>
 *   equal:ali,sara,reza            everyone owes the same
 *   shares:ali=1,sara=2            owed in proportion to the weights
 *   exact:ali=20.00,sara=30.00     the amounts are given directly
 * </pre>
 */
public final class Split {
    private static final Pattern WEIGHT = Pattern.compile("[0-9]+");

    private Split() {}

    public static Map<String, Long> parse(String spec, long amount) {
        int colon = spec.indexOf(':');
        if (colon < 0 || colon == spec.length() - 1) {
            throw new IllegalArgumentException("invalid split: " + spec);
        }
        String kind = spec.substring(0, colon);
        String rest = spec.substring(colon + 1);
        return switch (kind) {
            case "equal" -> equal(amount, Arrays.asList(rest.split(",", -1)));
            case "shares" -> byShares(amount, pairs(rest, Split::parseWeight));
            case "exact" -> exact(amount, pairs(rest, Money::parse));
            default -> throw new IllegalArgumentException("unknown split kind: " + kind);
        };
    }

    public static Map<String, Long> equal(long amount, List<String> names) {
        long each = amount / names.size();
        Map<String, Long> shares = new LinkedHashMap<>();
        for (String name : names) {
            shares.put(name, each);
        }
        return shares;
    }

    public static Map<String, Long> byShares(long amount, Map<String, Long> weights) {
        long total = weights.values().stream().mapToLong(Long::longValue).sum();
        Map<String, Long> shares = new LinkedHashMap<>();
        weights.forEach((name, weight) -> shares.put(name, amount * weight / total));
        return shares;
    }

    public static Map<String, Long> exact(long amount, Map<String, Long> parts) {
        long total = parts.values().stream().mapToLong(Long::longValue).sum();
        if (total > amount) {
            throw new IllegalArgumentException(
                    "exact parts add up to " + Money.format(total) + ", more than " + Money.format(amount));
        }
        return new LinkedHashMap<>(parts);
    }

    private static Map<String, Long> pairs(String text, ToLongFunction<String> parse) {
        Map<String, Long> pairs = new LinkedHashMap<>();
        for (String item : text.split(",", -1)) {
            int eq = item.indexOf('=');
            if (eq <= 0 || eq == item.length() - 1) {
                throw new IllegalArgumentException("expected name=value, got: " + item);
            }
            pairs.put(item.substring(0, eq), parse.applyAsLong(item.substring(eq + 1)));
        }
        return pairs;
    }

    private static long parseWeight(String text) {
        if (!WEIGHT.matcher(text).matches() || Long.parseLong(text) == 0) {
            throw new IllegalArgumentException("invalid weight: " + text);
        }
        return Long.parseLong(text);
    }
}
