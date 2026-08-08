# Questions

These are questions that you can ask before starting the hands-on interview to make sure
you are on the same page as the interviewee.

## How to Use This List

This is a **bank, not a script.** Reading it top to bottom produces a quiz, and a quiz measures
what someone memorised — not whether they can do the job. Pick a handful of entries that target
the signal you actually need, and turn each one into a conversation.

Four rules make the difference between an oral round that produces evidence and one that produces
a vague good feeling:

1. **Name the signal before you pick the question.** Write down what you are trying to learn —
   "can they reason about failure under concurrency?" — and then choose the entry that gets you
   there. If you cannot name it, skip the question.
2. **Prefer a scenario over a definition.** Most entries below are written as topics or as recall
   prompts. Recall has a binary answer and a thirty-second lifespan; a scenario has depth you can
   keep pulling on. See [Turning a Recall Prompt Into a Scenario](#turning-a-recall-prompt-into-a-scenario).
3. **Layer every question.** Have two or three follow-ups ready before you ask the first one.
   Depth is where the difference between candidates shows up — the opening question just gets
   everyone onto the same page.
4. **Do not play guess-what-I'm-thinking.** If you catch yourself saying "no, I meant something
   else," the question is underspecified. Restate it with the constraint you had in mind.

Record what the candidate said, not how it felt. `"could not explain why the index wasn't used,
even after a hint"` is evidence. `"database knowledge was weak"` is not.

## Levels

Same three levels as [`problems/`](../problems/), applied to oral questions. The level describes
which bar the question is calibrated against, not how obscure it is.

### Level 1 — Frontend and SRE

Everyday working knowledge. The candidate should have hit these situations in normal work, and
the answer should come from experience rather than study.

- **Budget:** 10–15 minutes across two or three topics.
- **Looking for:** concrete recall of things they have actually done, and honesty about what they
  have not. "I have never used `git bisect`, but I would go about it this way" is a good answer.
- **Do not** treat a missing term as a missing skill. Ask them to describe the situation instead.

### Level 2 — Backend engineers

Mechanism, not vocabulary. The candidate should be able to explain *why* something behaves the
way it does, and what breaks when it does not.

- **Budget:** 15–25 minutes on two topics, gone into properly.
- **Looking for:** a correct mental model that survives a follow-up. Push once past the first
  correct answer — the second question is the one that discriminates.
- **Do not** accept a textbook definition as the end of the answer. Ask for a case where it mattered.

### Level 3 — Seniors

Judgement under conflicting constraints. There is no single right answer, and the interesting part
is which trade-off they pick and whether they can defend it.

- **Budget:** 25–40 minutes, usually one scenario developed in depth.
- **Looking for:** naming the trade-off unprompted, changing their answer when you change the
  constraints, and saying what they would need to measure before deciding.
- **Do not** score these as right or wrong. Score the reasoning.

## Topics by Level

Most topics contain questions at more than one depth; this is where each section's centre of
gravity sits.

| Level | Topics |
|---|---|
| **1** | [Git](#git) · [Operating Systems](#operating-systems) · [Networking](#networking) · [Docker](#docker) · [Time Complexity](#time-complexity) · [Linked List](#linked-list) · [Sort Algorithms](#sort-algorithms) |
| **2** | [Golang](#golang) · [Python/Django](#pythondjango) · [Database](#database) · [Sessions](#sessions) · [Kubernetes](#kubernetes) · [SOLID](#solid) · [CI/CD](#cicd) · [Greedy Algorithms](#greedy-algorithms) |
| **3** | [System Design](#system-design) · [Cloud Native Design](#cloud-native-design) · [Design and Project structure](#design-and-project-structure) · [ML](#ml) · [Soft skills, Teamwork and Managerial](#soft-skills-teamwork-and-managerial) |

## Turning a Recall Prompt Into a Scenario

Several entries below are bare topic labels. They are useful as *reminders of what to cover*, but
asked verbatim they measure vocabulary. The rewrite is mechanical: take the fact, and build the
situation where not knowing it hurts.

| Instead of | Ask |
|---|---|
| "Foreign Key / Primary Key" | "Two services write to the same table and you start seeing orphaned rows. Walk me through how you would find the cause and stop it happening again." |
| "Container vs Virtual Machine" | "A container works on your laptop and crashes in production with the same image. What are the possible causes, and in what order would you check them?" |
| "Process vs Threads" | "A service's memory grows steadily until it is killed. How would you tell whether it is a leak, a thread pool, or the allocator?" |
| "What is Dependency Injection?" | "This class constructs its own database client. Show me what breaks when you try to test it, and what you would change." |
| "Do you understand the differences between Merge and Rebase?" | "You have been on a branch for two weeks and `main` has moved a long way. Talk me through how you get up to date, and what you would not do if others had pulled your branch." |

The pattern: the recall version has one right answer and ends; the scenario version has several
defensible answers and a follow-up for each.

## Known Weak Spots in This List

Flagged so you do not reach for them by accident, per the anti-patterns in
[`docs/designing-interview-questions.typ`](../docs/designing-interview-questions.typ):

- **Bare labels** (`Greedy Algorithms`, `Foreign Key`, `Primary Key`, `SOLID`, `Context`,
  `Embedding`) are section markers, not questions. Rewrite them into a scenario before asking.
- **Trivia with a lookup answer** (`free -m` vs `vmstat`, the exact `htop` state letters,
  the difference between `. ~/file` and `source ~/file`) tests memory for things everyone
  searches for. Keep them as warm-ups at most, and never let them affect the rating.
- **Questions with a hidden expected answer** (`Which sort has the best order among the
  comparison sorts?`) invite guess-what-I'm-thinking. Ask *why* the bound exists instead.
- **The `gossip` project question** assumes the candidate has worked on that specific
  repository. Generalise it to a project of *theirs* before asking.

## Git

**Level 1**

These questions are here to review the candidate's knowledge of `git`.
Believe it or not, there are many developers who don't know how to use `git` or Git Flow.
It is better to ask these questions in a storytelling manner. Describe the situation
and then ask how the candidate would solve it.

- Where have you used git?
- Do you understand the differences between Merge and Rebase?
- Are you familiar with Git flow?
- Have you used `git stash`?
- Have you used `git cherry-pick`?
- Have you used `git add -p ...`?
- Have you used `git bisect`?
- What are the differences between Git and GitHub?
- What process is an alternative to merging?
- How would you revert a commit that has already been pushed and made public?
- Can you remember some of your most used `git` commands?

## Algorithm

**Level 1–2**

### Time Complexity

- Definition of Time Complexity

### Linked List

- Differences between linked lists and arrays
- Time complexity for accessing an element
- An `ArrayList`, or dynamically resizing array, allows you to have the benefits of an array while offering flexibility in size.
  How do they achieve this? Consider adding `n` numbers into a `ArrayList`, what is the time complexity?

### Sort Algorithms

- Do you know any sorting algorithm that has O(n)?
- Which sort has the best order among the comparison sorts?
- What is the difference between Merge Sort and Quick Sort?

### Greedy Algorithms

## Sessions

**Level 2**

- How does session management work in a web application, and what are the different approaches to maintaining session state?

## Operating Systems

**Level 1**

- Process vs Threads
- Experience with multithreaded application programming
- How can you get the list of processes on Linux (`ps`)?
- Are you familiar with `grep`?
- What are the process states on Linux? (If the candidate does not know the process states' names, the questioner can describe at least 3)
  - Ready
  - Running
  - Blocked or wait
  - Terminated or Completed
  - Zombie
- What are the process states on Linux based on `htop`?
  - `S` for sleeping
  - `I` for idle (longer inactivity than sleeping on platforms that distinguish)
  - `R` for running
  - `D` for disk sleep (uninterruptible)
  - `Z` for zombie (waiting for parent to read its exit status)
  - `T` for traced or suspended (e.g by SIGTSTP)
  - `W` for paging
- Which command would you use to check how much memory is being used by Linux?
  - `free -m`
  - `vmstat`
  - `top`
  - `htop`
  - `cat /proc/meminfo`. (+)
- What is the difference between `. ~/file` and `~/file`
  - `./test.sh` runs `test.sh` as a separate program. It may happen to be a bash script,
    if the file `test.sh` starts with `#!/bin/bash`. But it could be something else altogether.
  - `. ./test.sh` executes the code of the file `test.sh` inside the running instance of bash.
    It works as if the content of file `test.sh` had been included textually instead of the `. ./test.sh` line.
    (Almost: there are a few details that differ, such as the value of `$BASH_LINENO`, and the behavior of the return built-in.)
- What is the difference between `. ~/file` and `source ~/file`
- Are you familiar with `systemd`?
- Why you cannot write on a disk that has the required spaces?
  - Because of the `inode` runs out

## Networking

**Level 1**

- When I type a URL on my laptop, can you tell me what my computer does?
  - The browser looks up the IP address of the server hosting the website. Your browser checks its own cache,
    the operating system cache, a local network cache at your router, and a DNS server cache on your corporate network
    or at your internet service provider (ISP).
    If the browser cannot find the IP address in any of those cache layers,
    the DNS server on your corporate network or at your ISP does a recursive DNS lookup.
    A recursive DNS lookup asks multiple DNS servers around the Internet,
    which in turn ask more DNS servers for the DNS record until it is found.
  - Browser initiates TCP connection with the server: Packets from a client browser request get routed through the
    router to find the server with the IP address to connect to. Instead, many sites use a content delivery network,
    or CDN, to cache static and dynamic content closer to the browser. Once the browser finds the server on the Internet,
    it establishes a TCP connection with the server and if HTTPS is being used,
    a TLS handshake takes place to secure the communication.
  - Browser sends the HTTP request to the server.
  - The server processes request and sends back a response.
  - Browser renders the content: As the browser is parsing and rendering the HTML,
    it is making additional requests to
    get JavaScript, CSS, images, and data. It can do much of this in parallel.
- How you can find the IP address of the server using its name?
- What are the differences between TCP and UDP?
- Flow Control vs Congestion Control
- How does a PHP request flow work (How does Common Gateway Interface (CGI) works)?
- Is there any restriction on the number of TCP connections for a system?
- Can you explain the process of TCP three-way handshake and the significance of each step?
  - SYN (Synchronize)
  - SYN-ACK (Synchronize-Acknowledge)
  - ACK (Acknowledgment)

## Python/Django

**Level 2**

- Have you had any experience with optimizing Django/Python projects?
  - <https://github.com/jazzband/django-silk>
  - Indexing
  - Pagination
- Are you familiar with Django Signals?
  - Django includes a **signal dispatcher** which helps decoupled applications get notified when actions
    occur elsewhere in the framework.
    In a nutshell, signals allow certain senders to notify a set of receivers that some action has taken place.
    They're especially useful when many pieces of code may be interested in the same events.
- Does following code have any issue?

  ```python
  async function_name(response):
    data = await response.json()
    # insert in db
  ```

- What is async programming?
- What is the difference between async programming and multithreaded programming?
- What is the difference between FastAPI and Flask?
- When can't you use async programming? When the implementation of the code is pure python (GIL)
- What is the difference between the two following codes?

```python
import asyncio

await asyncio.sleep(10)
```

```python
import time

time.sleep(10)
```

## Golang

**Level 2**

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
- Buffered/Un-buffered Channels
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

**Level 2**

- Foreign Key
- Primary Key
- NoSQL vs SQL

## Kubernetes

**Level 2**

- You have incidents in which your pod crashes randomly some minutes after its startup. What do you do about it? How you find out the problem?
- Did you write a Kubernetes manifest?
- Why we need _service_ for accessing to Kubernetes pods?
- Can we use pod's IP address for getting access to it?
- What are the differences between readiness and liveness probes?
- Do you know `helm`, `kustomize`, etc.?
- Can you explain the distinctions between statefulset and deployment?

## Docker

**Level 1**

- Container vs Virtual Machine
- How we can improve the following Dockerfile?

```dockerfile
# by adding the following file, you can use docker cache better
# COPY Pipfile Pipfile.lock ./

RUN pipenv install --dev --system --deploy

COPY . .
```

## SOLID

**Level 2**

- **S**: Single Responsibility Principle (known as SRP)
- **O**: Open/Closed Principle
- **L**: Liskov's Substitution Principle
- **I**: Interface Segregation Principle
- **D**: Dependency Inversion Principle

## Cloud Native Design

**Level 3**

- How do you handle a crashed loop application on Kubernetes?
- How do you monitor an application?
  - Metrics (Telemetry)
  - Logs
  - Tracing

## System Design

**Level 3**

- What do you know about deployment?
- Let's discuss one of these scenarios in detail
  - Event Delivery based on `MQTT`, `HTTP`, etc.
  - URL Shortener which contains
    - _Redis_
    - _Database Replication/Sharding_
    - _HAProxy_
    - ...
  - Voting System that introduces the **CAP** theorem
- What is Dependency Injection?

### Deployment Strategy

- What is Blue Green Deployment?
- What is A/B testing?
- How is A/B testing any different from Blue Green Deployment?
  1. Chat Service
  2. Routing Engine
- Have you used A/B testing before? What did you do?
  1. White List
  2. ID
  3. Rule Engine

## CI/CD

**Level 2**

If he/she used CI/CD:

- What are the benefits of CI/CD to deploy your code?
- Which one of tools (GitLab/Jenkins/Bitbucket) did you use?

If he/she did not use CI/CD:

- How did you deploy if? (Using Ansible or Puppet don't have negative point)
- Way did not use CI/CD?

## ML

**Level 3**

Our GPU doesn't have sufficient memory to load our model into it, what is your solution?

- Reduce Model Size or Use a Different Model Architecture:
  - Can you use a smaller pre-trained model?
  - Can you choose a more lightweight model architecture that is specifically designed for your task? For example, if you are working with
    deep learning, can you use MobileNet or SqueezeNet, which are designed to be more memory-efficient for tasks like image classification?
- Quantization: It can significantly reduce the memory footprint of a model by converting model weights to lower precision (e.g., from 32-bit floating point to 16-bit fixed point). Tools like TensorFlow's "tf.lite" or PyTorch's quantization modules can help with this.
- Use Mixed Precision Training
- Distributed Training and Model Parallelism
- Gradient Accumulation: This can allow you to use a larger batch size without increasing memory requirements.
- Prune or Sparsify the Model
- Optimize Your Code (Works for training): Free up GPU memory as soon as it's no longer needed.

- [Find S Algorithm](https://www.geeksforgeeks.org/ml-find-s-algorithm/)
- Dimensionality reduction reduces collinearity
- [Probably approximately correct (PAC) learning](https://www.baeldung.com/cs/probably-aproximately-correct)

Approximate Nearest Neighbor

There are so many ANN approaches, one of them is ANNOY

- ANNOY (Approximate Nearest Neighbor Oh Yeah):
  ![ANN](https://github.com/1995parham-teaching/interviews/assets/36500888/ac2a334d-5769-4fc0-a7a2-8fa87c2875d3)

## Soft skills, Teamwork and Managerial

**Level 3**

How do you prevent unwanted deployment of a new joiner?

1. I restrict production access for new joiners.
2. I disable all forms of automatic deployments for production, such as automatic synchronization in Argo CD.
3. I verify the image tag on GitLab to enable easy reversion if necessary.
4. I'll make main branch protected, so it will need approval to merge other branches with it.
5. I ensure monitoring alerts are set up to receive notifications in case of any code malfunctions.
6. I will implement a system integrated across all projects to send notifications to a deployment group whenever a
   project is deployed.

How do you make sure the on-calls can handle the incidents of your project?

1. Maintain an Incident Response Handbook.
2. Post-Incident Reviews
3. Mentorship and Shadowing
4. Simulation and Drills
5. All projects should implement effective monitoring systems.
6. Automation: for example having back up for the services we can like EMQ or having circuit breaker in projects.
7. Increase bus factor of projects.
