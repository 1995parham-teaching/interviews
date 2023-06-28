# Questions

## Git

These questions are here to review the level of candidate knowledge on `git`.
Trust me, there are many developers that don't know how to use `git` or Git Flow.
It is better to ask these questions in a story-telling manner. Describe the situation
and then ask how do you solve it.

- Where did you use git?
- Do you differences between Merge vs Rebase?
- Do you know Git flow?
- Did you use `git stash`?
- Did you use `git cherry-pick`?
- Did you use `git add -p ...`?
- Differences between Git and GitHub
- What process is an alternative to merging?
- How do you revert a commit that has already been pushed and made public?
- Do you remember some of your most used `git` command?

## Algorithm

- Do you any sorting algorithm that has O(n)?

## Python/Django

- Do you know about Django Signals?
  > Django includes a “signal dispatcher” which helps decoupled applications get notified when actions
  > occur elsewhere in the framework.
  > In a nutshell, signals allow certain senders to notify a set of receivers that some action has taken place.
  > They’re especially useful when many pieces of code may be interested in the same events.

## Golang

- Is there anything wrong with having more hierarchy for the Go package?
- Can you discuss the ways you have for creating a URL shortener service?
- Is there any restriction on the number of TCP connections for a system?

- Talk about [gossip](https://github.com/elahe-dastan/gossip) project and try to tell the following points:

  - Challenges
  - Mistakes/Failure
  - Enjoyed
  - What you'd do differently

- An ArrayList, or dynamically resizing array, allows you to have the benefits of an array while offering flexibility in size.
  How do they do this? consider we want to add n number into ArrayList what is the time complexity?

- Can you explain the following cases in Golang:

```go
var ch1 chan int
ch1 <- 1 // write on a nil channel
<-ch1 // read from a nil channel

ch2 := make(chan int)
close(ch2)
ch2 <- 1 // write on a closed channel
<-ch2 // read from a closed channel
```

[Answer](https://stackoverflow.com/questions/39015602/how-does-a-non-initialized-channel-behave)

- Explain panics in Golang. Can you mention some of these cases?

- Explain the result of the following code. Is there any issue?

```go
package main

import "fmt"

func main() {
        s := []int{1, 2, 3, 4}

        change(s)

        fmt.Println(s)
}

func change(s []int) {
        t := make([]int, len(s))

        copy(t, s)

        for i := range t {
                t[i]++
        }

        s = t
}
```

- Explain Reader-Writer problem.
- Explain a mutex implementation in Go.
- Describe the error handling procedure in Go and error wrapping.
- What is the context? How we can use it to cancel the long-run processing?
- Is there any difference between `array` and `slice` in Golang?
- You have incidents in which your pod crashes randomly some minutes after its startup. What do you do about it? How you find out the problem?

## Database

- Foreign Key
- Primary Key
- NoSQL vs SQL

## Docker

- Container vs Virtual Machine

## SOLID

- **S**: Single Responsibility Principle (known as SRP)
- **O**: Open/Closed Principle
- **L**: Liskov’s Substitution Principle
- **I**: Interface Segregation Principle
- **D**: Dependency Inversion Principle
