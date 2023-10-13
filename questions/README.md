# Questions

These are questions that you can ask before starting the hands-on interview to make sure
you are on a same page with the interviewee.

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
- Did you use `git bisect`?
- Differences between Git and GitHub
- What process is an alternative to merging?
- How do you revert a commit that has already been pushed and made public?
- Do you remember some of your most used `git` command?

## Algorithm

### Time Complexity

- Time Complexity Definition

### Linked List

- Differences between linked list and arrays
- Time complexity for accessing an element
- An `ArrayList`, or dynamically resizing array, allows you to have the benefits of an array while offering flexibility in size.
  How do they do this? Consider we want to add `n` number into a `ArrayList`, what is the time complexity?

### Sort Algorithms

- Do you any sorting algorithm that has O(n)?
- Which sort has the best order among the comparison sorts?
- What is the difference between Merge Sort and Quick Sort?

### Greedy Algorithms

## Sessions

- How does session management work in a web application, and what are the different approaches to maintaining session state?

## Operating Systems

- Process vs Threads
- Multi-thread application programming experience
- How we can get the processes list on Linux (`ps`)
- Do you know `grep`?
- What are the process states on Linux? (If candidate does not know about the process states name, questioner can read below states and wait to describe at least 3)
  - Ready
  - Running
  - Blocked or wait
  - Terminated or Completed
  - Zombie
- What command would you use to check how much memory is being used by Linux?
  - `free -m`
  - `vmstat`
  - `top`
  - `htop`
  - `cat /proc/meminfo`. (+)
- What is the difference between `. ~/file` and `~/file`
  - `./test.sh` runs test.sh as a separate program. It may happen to be a bash script, if the file test.sh starts with `#!/bin/bash`. But it could be something else altogether.
  - `. ./test.sh` executes the code of the file `test.sh` inside the running instance of bash. It works as if the content file `test.sh` had been included textually instead of the `. ./test.sh` line.
    (Almost: there are a few details that differ, such as the value of `$BASH_LINENO`, and the behavior of the return built-in.)
- What is the difference between `. ~/file` and `source ~/file`

## Networking

- When I type a URL in my laptop can you tell me what my computer does?
- What are the differences between TCP and UDP?
- Flow Control vs Congestion Control
- How does a PHP request flow work (How does Common Gateway Interface (CGI) works)?
- Is there any restriction on the number of TCP connections for a system?
- Can you explain the process of TCP three-way handshake and the significance of each step?
  - SYN (Synchronize)
  - SYN-ACK (Synchronize-Acknowledge)
  - ACK (Acknowledgment)

## Python/Django

- Did you have any experience with optimizing Django/Python project?
  - <https://github.com/jazzband/django-silk>
  - Indexing
  - Pagination
- Do you know about Django Signals?
  - Django includes a “signal dispatcher” which helps decoupled applications get notified when actions
    occur elsewhere in the framework.
    In a nutshell, signals allow certain senders to notify a set of receivers that some action has taken place.
    They’re especially useful when many pieces of code may be interested in the same events.

## Golang

- What are the differences between value types and reference types in Golang?
  - Value types
    - Store their data directly in the memory where the variable is allocated.
    - Primitive types like `int`, `float32`, `float64`, `bool`, `struct`, and `array` are **value** types in Golang.
  - Reference Types
    - Store a reference (i.e., memory address) to the actual data, which is stored on the heap.
- What is the purpose of `defer` statements in Golang?
  - Execution Order: Multiple defer statements within a function are executed in a last-in, first-out (LIFO) order. The deferred function calls are placed on a stack, and as the surrounding function exits, the calls are executed in reverse order.

### Design and Project structure

- Is there anything wrong with having more hierarchy for the Go package?
- Can you discuss the ways you have for creating a URL shortener service?
- Talk about [gossip](https://github.com/elahe-dastan/gossip) project and try to tell the following points:
  - Challenges
  - Mistakes/Failure
  - Enjoyed
  - What you'd do differently

### Arrays vs Slices

- A `ArrayList`, or dynamically resizing array, allows you to have the benefits of an array while offering flexibility in size.
  How do they do this? Consider we want to add n number into `ArrayList` what is the time complexity?

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

- Is there any difference between `array` and `slice` in Golang?

### Context

- What is the context? How we can use it to cancel the long-run processing?

### Goroutine

- How does Goroutine differ from a regular thread?

### Channels and Synchronization

- Did you use channels? Where did you use them?
- Buffered/Unbuffered Channels
- `select`
- Sync Package (Mutex and Semaphore, WaitGroup)
- Solve Reader-Writer problem with channels
- Explain mutex implementation with channels

### Embedding

```go
type Student struct {
    Person
}

type Person struct {
    Name string
    Age int
}
```

### Interfaces

- How they are different from Java interfaces?

### Testing

- Have you ever written tests for you Go projects?
- Have you ever used _mock_ in your projects?

### Empty Structure and Why?

```go
type Empty struct {}
```

### Errors

- Describe the error handling procedure in Go and error wrapping.
- Explain `panic` in Golang. Can you mention some of these cases?
  - Out of Bounds Panics
  - Nil Receivers

## Database

- Foreign Key
- Primary Key
- NoSQL vs SQL

## Kubernetes

- You have incidents in which your pod crashes randomly some minutes after its startup. What do you do about it? How you find out the problem?
- Did you write a Kubernetes manifest?
- Why we need _service_ for accessing to Kubernetes pods?
- Can we use pod's IP address for getting access to it?
- What are the differences between readiness and liveness probes?
- Do you know `helm`, `kustomize`, etc.?

## Docker

- Container vs Virtual Machine
- How we can improve the following Dockerfile?

```dockerfile
# by adding the following file, you can use docker cache better
# COPY Pipfile Pipfile.lock ./

RUN pipenv install --dev --system --deploy

COPY . .
```

## SOLID

- **S**: Single Responsibility Principle (known as SRP)
- **O**: Open/Closed Principle
- **L**: Liskov’s Substitution Principle
- **I**: Interface Segregation Principle
- **D**: Dependency Inversion Principle

## Cloud Native Design

- How do you handle a crashed loop application on Kubernetes?
- How do you monitor an application?
  - Metrics (Telemetry)
  - Logs
  - Tracing

## System Design

- What do you know about deployment?
- Let's discuss one of these scenarios in detail
  - Event Delivery based on `MQTT`, `HTTP`, etc.
  - URL Shortener which contains
    - _Redis_
    - _Database Replication/Sharding_
    - _HAProxy_
    - ...
  - Voting System that introduces the **CAP** theorem

## CI/CD

If he/she used CI/CD:

- What are the benefits of CI/CD to deploy your code?
- Which one of tools (Gitlab/Jenkins/Bitbucket) did you use?

If he/she did not use CI/CD:

- How did you deploy if? (Using Ansible or Puppet don’t have negative point)
- Way did not use CI/CD?
