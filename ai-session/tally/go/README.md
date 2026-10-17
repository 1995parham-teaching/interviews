# Tally

Tally keeps a shared expense ledger for a group of people: a trip, a flat, a dinner. People pay for things on behalf of the group, Tally works out who owes whom, and suggests the payments that settle everything.

There are no dependencies beyond the standard library. Go 1.22 or later.

## Running it

```
go run ./cmd/tally examples/trip.txt
```

Without a file, commands are read from standard input, one per line.

## Commands

```
member <name>
expense <payer> <amount> <split> <description>
pay <sender> <receiver> <amount>
delete <expense-id>
balances
settle
```

Blank lines and lines starting with `#` are ignored. A failing command prints `error: line N: ...` and the run continues; the exit code is 1 if any command failed.

- `member` adds a person to the group. Names are single words.
- `expense` records that `payer` paid `amount` for the group. Amounts are decimal, such as `12.50` or `7`. The split says who the expense was for:
  - `equal:ali,sara,reza` splits it equally.
  - `shares:ali=1,sara=2` splits it in proportion to the weights.
  - `exact:ali=20.00,sara=30.00` gives each person's part directly.

  The payer does not have to be part of the split. Each expense gets an id, printed as `#N`.
- `pay` records that `sender` gave `receiver` money directly, usually to pay back a debt.
- `delete` removes an expense by id, as if it had never been recorded.
- `balances` prints each member's balance. Positive means the group owes them money; negative means they owe the group. The balances of a group always add up to zero.
- `settle` prints a list of payments that would bring every balance back to zero.

## Example

```
$ go run ./cmd/tally examples/trip.txt
#1 hotel: ali paid 90.00
#2 dinner: sara paid 24.00
#3 train tickets: reza paid 18.00
sara paid ali 10.00
ali   +36.00
reza  -24.00
sara  -12.00
reza -> ali 24.00
sara -> ali 12.00
```

## Tests

```
go test ./...
```
