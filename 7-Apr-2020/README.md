# Problem.
We have n amount of money and our country have the following coins:

coin-1
coin-5
coin-7
coin-10

we want to have this money with minimum number of coins. what is the minimum?
for eaxmple:
- 2 = 2 * coin-1
- 5 = 1 * coin-5
- 6 = 1 * coin-5 + 1 * coin-1

# Questions.

1. [hub](https://github.com/elahe-dastan/hub) Why do you have a request/response package? What do these packages do?
2. [hub](https://github.com/elahe-dastan/hub) Why do you use [Koanf](https://github.com/knadh/koanf)? What is wrong with [Viper](https://github.com/spf13/viper)?
3. [URLshortener](https://github.com/elahe-dastan/urlShortener) Why don't you add the URL shortener docker into its docker-compose?
4. Is there anything wrong with having more hierarchy for the Go package?
5. Can you discuss the ways you have for creating a URL shortener service?
6. Is there any restriction on the number of TCP connections for a system?
7. Is there any sorting algorithm that has O(n)?

# System Design.
You are the CTO of a Fashion Startup.
This fashion startup has the voting system to select the `Beshooest Person Ever`. This voting system has the following structure:

```

        +--------+
        |   HQ   |
        +--------+


+----+  +----+  +----+  +----+
| C1 |  | C2 |  | C3 |  | C4 |  ....
+----+  +----+  +----+  +----+

```

As a CTO can you propose a high-level design for this voting system? If only we want the results at the end of the voting period? what happens if we want the real-time result?

1. There is no issue in the network
2. There are issues in the network that happens randomly and may disturb your system.
3. People may fraud in election results by voting many times.
