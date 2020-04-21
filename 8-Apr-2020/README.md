# Questions.

1. Explain Reader-Writer problem.
2. Explain mutex implementation in Go.
3. Describe the error handling procedure in Go and error wrapping.
4. What is the context? How we can use it to cancel the long-run processing?

# Code Session.

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

if the player runs out of the score he or she will die.
