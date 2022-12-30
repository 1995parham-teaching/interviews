# Questions

1. Is there anything wrong with having more hierarchy for the Go package?
2. Can you discuss the ways you have for creating a URL shortener service?
3. Is there any restriction on the number of TCP connections for a system?
4. Is there any sorting algorithm that has O(n)?

5. Talk about [gossip](https://github.com/elahe-dastan/gossip) project and try to tell the following points:

   - Challenges
   - Mistakes/Failure
   - Enjoyed
   - What you'd do differently

6. An ArrayList, or dynamically resizing array, allows you to have the benefits of an array while offering flexibility in size.
   How do they do this? consider we want to add n number into ArrayList what is the time complexity?

7. Can you explain the following cases in Golang:

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

8. Explain panics in Golang. Can you mention some of these cases?

9. Do you know `docker`? How it works?

10. Explain the result of the following code. Is there any issue?

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

11. Explain Reader-Writer problem.
12. Explain mutex implementation in Go.
13. Describe the error handling procedure in Go and error wrapping.
14. What is the context? How we can use it to cancel the long-run processing?
15. Is there any difference between `array` and `slice` in Golang?
