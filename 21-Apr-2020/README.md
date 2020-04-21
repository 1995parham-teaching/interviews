# Questions.

1. Talk about [gossip](https://github.com/elahe-dastan/gossip) project and try to tell the following points:
    - Challenges
    - Mistakes/Failure
    - Enjoyed
    - What you'd do differently

2. An ArrayList, or dynamically resizing array, allows you to have the benefits of an array while offering flexibility in size.
How they do this? consider we want to add n number into ArrayList what is the time complexity?

3. Can you explain the following cases in Golang:

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

4. Explain panics in Golang. Can you mention some of these cases?

5. Do you know `docker`? How it works?

# Problems

1. Consider a m x n matrix that each of its rows is soreted. Write an efficient algorithm to find a given number into it.
