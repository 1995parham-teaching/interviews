import io
import unittest
from pathlib import Path

from tally.cli import run

EXAMPLE = Path(__file__).parent.parent / "examples" / "trip.txt"


class CliTest(unittest.TestCase):
    def test_example_trip(self):
        out = io.StringIO()
        with EXAMPLE.open(encoding="utf-8") as f:
            failures = run(f, out)
        self.assertEqual(failures, 0)
        self.assertEqual(
            out.getvalue(),
            "#1 hotel: ali paid 90.00\n"
            "#2 dinner: sara paid 24.00\n"
            "#3 train tickets: reza paid 18.00\n"
            "sara paid ali 10.00\n"
            "ali   +36.00\n"
            "reza  -24.00\n"
            "sara  -12.00\n"
            "reza -> ali 24.00\n"
            "sara -> ali 12.00\n",
        )

    def test_errors_do_not_stop_the_run(self):
        out = io.StringIO()
        failures = run(["member ali", "pay ali mina 5", "member sara", "balances"], out)
        self.assertEqual(failures, 1)
        self.assertEqual(
            out.getvalue(),
            "error: line 2: unknown member: mina\n"
            "ali   0.00\n"
            "sara  0.00\n",
        )


if __name__ == "__main__":
    unittest.main()
