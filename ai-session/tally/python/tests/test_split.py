import unittest

from tally.split import parse_split


class SplitTest(unittest.TestCase):
    def test_equal(self):
        self.assertEqual(
            parse_split("equal:ali,sara,reza", 9000),
            {"ali": 3000, "sara": 3000, "reza": 3000},
        )

    def test_shares(self):
        self.assertEqual(
            parse_split("shares:ali=1,sara=2", 3000),
            {"ali": 1000, "sara": 2000},
        )

    def test_exact(self):
        self.assertEqual(
            parse_split("exact:ali=20.00,sara=30.00", 5000),
            {"ali": 2000, "sara": 3000},
        )

    def test_exact_more_than_amount(self):
        with self.assertRaises(ValueError):
            parse_split("exact:ali=20.00,sara=40.00", 5000)

    def test_invalid(self):
        for spec in ["equal", "equal:", "halves:ali,sara", "shares:ali", "shares:ali=0"]:
            with self.subTest(spec=spec), self.assertRaises(ValueError):
                parse_split(spec, 1000)


if __name__ == "__main__":
    unittest.main()
