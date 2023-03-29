# System Design

## Beshooest Person Ever (A Distributed Voting System)

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

As a CTO can you propose a high-level design for this voting system?
If only we want the results at the end of the voting period?
What happens if we want the real-time result?

1. There is no issue in the network
2. There are issues in the network that happens randomly and may disturb your system.
3. People may use fraud in election results by voting many times.

The objective is talking about CAP theorem and make sure he/she understands it correctly.

## Number Masking

Snapp! wants to have a service for hiding the driver and passenger numbers to each other. How do you implement it?

The objective is talking about how do you handle the masking process life cycle, how do you store the user preferences, etc.

## Event Engine

Snapp! wants to deliver events into drivers and passengers.
Which network protocol(s) suites for this problem?
Event Delivery based on MQTT, HTTP, etc.

## URL Shortener

We want to design a system that reads a URL and returns a short version of it (this version should be more memorable than the original one).
It is better to use following technologies:

- Redis
- Database Replication/Sharding
- HAProxy

## How do you react about an incident?

How do you find the malfunctioning service in case that you have so many microservices,
and your response time is going up?

The objective is talking about Prometheus, Grafana and Jeager.
