package tally;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.io.BufferedReader;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.PrintStream;
import java.io.StringReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import org.junit.jupiter.api.Test;

class CliTest {
    @Test
    void exampleTrip() throws IOException {
        ByteArrayOutputStream bytes = new ByteArrayOutputStream();
        int failures;
        try (BufferedReader in = Files.newBufferedReader(Path.of("examples/trip.txt"), StandardCharsets.UTF_8)) {
            failures = Cli.run(in, new PrintStream(bytes, true, StandardCharsets.UTF_8));
        }
        assertEquals(0, failures);
        assertEquals(
                """
                #1 hotel: ali paid 90.00
                #2 dinner: sara paid 24.00
                #3 train tickets: reza paid 18.00
                sara paid ali 10.00
                ali   +36.00
                reza  -24.00
                sara  -12.00
                reza -> ali 24.00
                sara -> ali 12.00
                """,
                bytes.toString(StandardCharsets.UTF_8));
    }

    @Test
    void errorsDoNotStopTheRun() throws IOException {
        ByteArrayOutputStream bytes = new ByteArrayOutputStream();
        BufferedReader in = new BufferedReader(new StringReader("member ali\npay ali mina 5\nmember sara\nbalances\n"));
        int failures = Cli.run(in, new PrintStream(bytes, true, StandardCharsets.UTF_8));
        assertEquals(1, failures);
        assertEquals(
                """
                error: line 2: unknown member: mina
                ali   0.00
                sara  0.00
                """,
                bytes.toString(StandardCharsets.UTF_8));
    }
}
