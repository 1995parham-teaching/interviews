import unittest

from tally.money import format_amount, parse_amount


class ParseAmountTest(unittest.TestCase):
    def test_whole_and_cents(self):
        self.assertEqual(parse_amount("12.50"), 1250)
        self.assertEqual(parse_amount("7"), 700)
        self.assertEqual(parse_amount("0.05"), 5)

    def test_rejects_garbage(self):
        for text in ["", "abc", "-3", "1,5", "0"]:
            with self.subTest(text=text), self.assertRaises(ValueError):
                parse_amount(text)


class FormatAmountTest(unittest.TestCase):
    def test_format(self):
        self.assertEqual(format_amount(1250), "12.50")
        self.assertEqual(format_amount(5), "0.05")
        self.assertEqual(format_amount(-300), "-3.00")

    def test_signed(self):
        self.assertEqual(format_amount(4000, signed=True), "+40.00")
        self.assertEqual(format_amount(0, signed=True), "0.00")


if __name__ == "__main__":
    unittest.main()
