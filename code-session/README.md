# Code Session

## Record Appender

There is a data file as below. We want to read it and then insert it into database.

data.txt:

```
John, Doe, 0939 123 1234
Jane, Doe, 0399 123 1234
```

- API for retrieving data from database
- Improve performance of Insert phase to support thousands of records

P.S. [Solution](https://github.com/1995parham-teaching/record-appender)

## Hangman

Implement Hangman!

```
> Computer: _ _ _ _ _

< Player: P

> Computer _ _ _ _ _ (player lose some score)

< Player: E

> Computer E _ _ _ E

< Player: L

> Computer: E L _ _ E

< Player: H

> Computer E L _ H E

< Player A

> Computer E L A H E (player win)
```

If the player runs out of the score he or she will die.
