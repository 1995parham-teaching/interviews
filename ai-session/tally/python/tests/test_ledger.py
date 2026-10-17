import unittest

from tally import settle
from tally.ledger import Ledger, LedgerError
from tally.settle import Transfer


def trip() -> Ledger:
    ledger = Ledger()
    for name in ["ali", "sara", "reza"]:
        ledger.add_member(name)
    return ledger


class LedgerTest(unittest.TestCase):
    def test_expense_updates_balances(self):
        ledger = trip()
        ledger.add_expense("ali", 9000, {"ali": 3000, "sara": 3000, "reza": 3000}, "hotel")
        self.assertEqual(ledger.balances(), {"ali": 6000, "sara": -3000, "reza": -3000})

    def test_payment_updates_balances(self):
        ledger = trip()
        ledger.add_expense("ali", 9000, {"ali": 3000, "sara": 3000, "reza": 3000}, "hotel")
        ledger.record_payment("sara", "ali", 3000)
        self.assertEqual(ledger.balances(), {"ali": 3000, "sara": 0, "reza": -3000})

    def test_delete_reverts_expense(self):
        ledger = trip()
        ledger.add_expense("ali", 9000, {"ali": 3000, "sara": 3000, "reza": 3000}, "hotel")
        dinner = ledger.add_expense("sara", 2000, {"ali": 1000, "sara": 1000}, "dinner")
        ledger.delete_expense(dinner.id)
        self.assertEqual(ledger.balances(), {"ali": 6000, "sara": -3000, "reza": -3000})
        self.assertEqual([e.description for e in ledger.expenses()], ["hotel"])

    def test_unknown_member(self):
        ledger = trip()
        with self.assertRaises(LedgerError):
            ledger.add_expense("mina", 1000, {"ali": 1000}, "coffee")
        with self.assertRaises(LedgerError):
            ledger.record_payment("ali", "mina", 1000)


class SettleTest(unittest.TestCase):
    def test_plan(self):
        transfers = settle.plan({"ali": 3600, "sara": -1200, "reza": -2400})
        self.assertEqual(
            transfers,
            [Transfer("reza", "ali", 2400), Transfer("sara", "ali", 1200)],
        )

    def test_nothing_to_settle(self):
        self.assertEqual(settle.plan({"ali": 0, "sara": 0}), [])


if __name__ == "__main__":
    unittest.main()
