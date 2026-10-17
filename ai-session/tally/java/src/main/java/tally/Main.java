package tally;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;

public final class Main {
    private Main() {}

    public static void main(String[] args) throws IOException {
        if (args.length > 1) {
            System.err.println("usage: tally [file]");
            System.exit(2);
        }
        int failures;
        if (args.length == 1) {
            try (BufferedReader in = Files.newBufferedReader(Path.of(args[0]), StandardCharsets.UTF_8)) {
                failures = Cli.run(in, System.out);
            }
        } else {
            BufferedReader in = new BufferedReader(new InputStreamReader(System.in, StandardCharsets.UTF_8));
            failures = Cli.run(in, System.out);
        }
        System.exit(failures > 0 ? 1 : 0);
    }
}
