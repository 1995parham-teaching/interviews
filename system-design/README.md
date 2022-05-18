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
