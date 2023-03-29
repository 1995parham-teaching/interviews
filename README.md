# Interviews

## Introduction

This repository started as mocked interviews between me and [@elahe-dastan](https://github.com/elahe-dastan) then I decided to gather all together at our team organization.
Now this repository contains sample problems for all stages in technical interview.
Please fill issue in case of any problem with these questions.

## Before Interview

1. Only have these tabs open:

- Gmail: Maybe they want to share a Google document with you
- [GoByExample](https://gobyexample.com/): Maybe they let you review concepts in their live coding session,
  so having this at your hand to find out your challenges.
- [Golang](https://pkg.go.dev/)

2. Have Ghermezi up and running
3. It is ok that you talk to yourself but try to do dispose of everything you don't know :)
4. Don't have any stress, they are as senior as you :D

## Sessions

### 6 Apr 2020

- [Problem](./problems#Shuffle)

### 7 Apr 2020

- [Problem](./problems#Coins)
- [Problem](./system-design#bshooest-person-ever)

1. [hub](https://github.com/elahe-dastan/hub) Why do you have a request/response package? What do these packages do?
2. [hub](https://github.com/elahe-dastan/hub) Why do you use [Koanf](https://github.com/knadh/koanf)? What is wrong with [Viper](https://github.com/spf13/viper)?
3. [URLshortener](https://github.com/elahe-dastan/urlshortener) Why don't you add the URL shortener docker into its docker-compose?
4. Is there anything wrong with having more hierarchy for the Go package?
5. Can you discuss the ways you have for creating a URL shortener service?
6. Is there any restriction on the number of TCP connections for a system?
7. Is there any sorting algorithm that has O(n)?

### 8 Apr 2020

- Discuss Memorization and Dynamic programming. Can you implement the coin problem?
- Review on Reader-Writer Problem
- How to implement semaphore with channels?
- Discuss error wrapping and how it can be useful.
- What is the context? How we can use it to cancel the long-run processing?

- [Problem](./code-session#Hangman)

### 9 Apr 2020

- [Problem](./code-session#record-appender)
- [Problem](./problems#snappfood)
- [Problem](./problems#happy-number)

### 10 Apr 2020

Coding session:

Review on Golang concurrency concepts + HTTP APIs with SnappFood project.

Review on these points:

1. Your pipeline idea is great
2. Time-based writing to the database is better than count-based
3. Build strings with a string builder, but your way is OK too
4. more about variadic functions and how to convert slices to parameters and vice versa.
5. Thread Pools and Golang Go Routines Pipeline
6. Query Parameter in HTTP Request
7. Do you know what is happening when you enter `www.google.com` in your browser?
8. Auto-Increment ID in PostgreSQL
9. Create/Drop Tables in PostgreSQL
10. Insert/Select/Delete from Tables in PostgreSQL
11. Review on HTTP Status Codes is Good
12. Git Reset/Rebase
13. Docker

### 21 Apr 2022

- [Problem](./problems/merge)
