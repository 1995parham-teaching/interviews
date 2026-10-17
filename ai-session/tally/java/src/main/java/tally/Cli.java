package tally;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.PrintStream;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;

/** Reads commands, one per line, and runs them against a single ledger. */
public final class Cli {
    private static final String USAGE =
            """
            commands:
              member <name>
              expense <payer> <amount> <split> <description>
              pay <sender> <receiver> <amount>
              delete <expense-id>
              balances
              settle""";

    private static final Pattern ID = Pattern.compile("[0-9]+");

    private Cli() {}

    /** Runs every command and returns the number of commands that failed. */
    public static int run(BufferedReader in, PrintStream out) throws IOException {
        Ledger ledger = new Ledger();
        int failures = 0;
        int number = 0;
        String raw;
        while ((raw = in.readLine()) != null) {
            number++;
            String line = raw.strip();
            if (line.isEmpty() || line.startsWith("#")) {
                continue;
            }
            try {
                execute(ledger, line, out);
            } catch (LedgerException | IllegalArgumentException e) {
                out.println("error: line " + number + ": " + e.getMessage());
                failures++;
            }
        }
        return failures;
    }

    public static void execute(Ledger ledger, String line, PrintStream out) {
        String[] fields = line.split("\\s+");
        String command = fields[0];
        List<String> args = Arrays.asList(fields).subList(1, fields.length);
        switch (command) {
            case "member" -> {
                arity(args, 1, "member <name>");
                ledger.addMember(args.get(0));
            }
            case "expense" -> {
                if (args.size() < 4) {
                    throw new IllegalArgumentException("usage: expense <payer> <amount> <split> <description>");
                }
                long amount = Money.parse(args.get(1));
                Map<String, Long> shares = Split.parse(args.get(2), amount);
                String description = String.join(" ", args.subList(3, args.size()));
                Ledger.Expense expense = ledger.addExpense(args.get(0), amount, shares, description);
                out.printf("#%d %s: %s paid %s%n",
                        expense.id(), expense.description(), expense.payer(), Money.format(expense.amount()));
            }
            case "pay" -> {
                arity(args, 3, "pay <sender> <receiver> <amount>");
                Ledger.Payment payment = ledger.recordPayment(args.get(0), args.get(1), Money.parse(args.get(2)));
                out.printf("%s paid %s %s%n", payment.sender(), payment.receiver(), Money.format(payment.amount()));
            }
            case "delete" -> {
                arity(args, 1, "delete <expense-id>");
                if (!ID.matcher(args.get(0)).matches()) {
                    throw new IllegalArgumentException("invalid expense id: " + args.get(0));
                }
                Ledger.Expense expense = ledger.deleteExpense(Integer.parseInt(args.get(0)));
                out.printf("deleted #%d %s%n", expense.id(), expense.description());
            }
            case "balances" -> {
                arity(args, 0, "balances");
                printBalances(ledger.balances(), out);
            }
            case "settle" -> {
                arity(args, 0, "settle");
                List<Settle.Transfer> transfers = Settle.plan(ledger.balances());
                if (transfers.isEmpty()) {
                    out.println("nothing to settle");
                }
                for (Settle.Transfer t : transfers) {
                    out.printf("%s -> %s %s%n", t.sender(), t.receiver(), Money.format(t.amount()));
                }
            }
            default -> throw new IllegalArgumentException("unknown command: " + command + "\n" + USAGE);
        }
    }

    private static void printBalances(Map<String, Long> balances, PrintStream out) {
        int width = balances.keySet().stream().mapToInt(String::length).max().orElse(0);
        balances.keySet().stream().sorted().forEach(name ->
                out.printf("%-" + width + "s  %s%n", name, Money.formatSigned(balances.get(name))));
    }

    private static void arity(List<String> args, int n, String usage) {
        if (args.size() != n) {
            throw new IllegalArgumentException("usage: " + usage);
        }
    }
}
