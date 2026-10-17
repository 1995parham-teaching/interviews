"""Work out who should pay whom so that every balance goes back to zero."""

from dataclasses import dataclass


@dataclass(frozen=True)
class Transfer:
    sender: str
    receiver: str
    amount: int


def plan(balances: dict[str, int]) -> list[Transfer]:
    """Repeatedly let the member who owes the most pay the member who is owed the most.

    Ties are broken by name so that the plan is the same on every run.
    """
    transfers: list[Transfer] = []
    while balances:
        names = sorted(balances)
        creditor = max(names, key=lambda name: balances[name])
        debtor = min(names, key=lambda name: balances[name])
        if balances[creditor] <= 0 or balances[debtor] >= 0:
            break
        amount = min(balances[creditor], -balances[debtor])
        balances[creditor] -= amount
        balances[debtor] += amount
        transfers.append(Transfer(debtor, creditor, amount))
    return transfers
