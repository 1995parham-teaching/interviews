# AI Session

A code session where the candidate is expected to use an AI assistant, on a real codebase they have never seen. The question is not whether they finish. It is how they work with the assistant inside unfamiliar code, and whether they notice when the assistant is wrong.

This is a different instrument from [`code-session/`](../code-session/). A code session measures how someone builds from a blank page. An AI session measures how someone reads, verifies and steers when most of the typing is done by a tool. The two produce different evidence, and a loop that uses both learns more than a loop that uses either twice.

## Format

- **Individual, 120 minutes.** Four timeboxed tasks plus a stretch, then a 15 minute debrief.
- **Bring your own AI.** Any tool is allowed: Claude Code, Cursor, Copilot, a chat window. Record which one, and whether it can run commands by itself. A chat window forces the candidate to be the hands; an agentic tool shifts the question to whether they read what it ran.
- **Observed.** One interviewer with the notes sheet, plus screen recording. If the session is unobserved, the AI transcript and the git history carry most of the signal.
- **Commit per task.** The candidate tags each task's end state (`task-0`, `task-1`, ...). Partial work is fine and is itself evidence.
- **Tasks are revealed one at a time.** Later tasks touch bugs that earlier tasks may have exposed; revealing them early spoils the signal.
- **The environment is prepared.** Docker images pulled, dependencies installed. Setup friction is noise, not signal.

## Scoring

Score in the four competencies of the evaluation document, on its scale — Weak (ضعیف) · Borderline (مرزی) · Strong (قوی) · Not assessed (سنجیده نشد). The AI changes what counts as evidence for each competency, not the competencies themselves.

| Competency | What Weak looks like | What Strong looks like |
| --- | --- | --- |
| Data structures, algorithms and problem solving | Patches where the symptom appears; cannot say where state lives or what the metaclass does at the debrief | Finds the origin of each bug, checks what else it affected, draws the data flow unaided |
| Coding | Lets the AI rewrite the metaclass, add a framework, or move to FastAPI; tests pass before the fix | Keeps the diff to the task, writes a failing test first, follows the repo's existing patterns |
| Effective communication | Prompts like "agents is empty, why"; accepts the AI's explanation verbatim; API contract decided silently by the AI | Prompts carry the traceback, the file and the expectation; challenges a wrong explanation with evidence; names a design decision and asks or states an assumption |
| Speed and getting to a result | Applies each AI suggestion in turn until something works; never calls the endpoint or runs the tests | Reproduces before, verifies after, reads the diff, and still reaches working code in the time box |

Also record three plain counts per task: files opened by the candidate, files opened by the AI, and AI claims accepted without a check. These separate "reads code" from "reads about code" better than any score.

## MiA Cold Start

**Codebase:** [I1820/mia](https://github.com/I1820/mia) at commit `f0bd540`, a small IoT middleware. Agents publish device data over MQTT; users read and control devices over HTTP. About 1.9k lines of Python across 71 files, six tests, all green.

**Why this codebase.** It is small enough to read in an hour but built on Python features most people only half-remember: a metaclass that registers every `Thing` subclass by name, descriptors with side effects (assigning to a `Setting` publishes an MQTT message, assigning to a `State` writes to Mongo), singleton controllers wired by class attributes from `main.py`, and lazy `importlib` loading of thing types by module name. An AI assistant describes all of this fluently and gets details wrong. It also has a real runtime, so the candidate can verify claims end to end instead of only reading. Whether they choose to is the main thing we observe.

**Latent bugs.** The current `main` branch carries five bugs that CI does not catch. They were reproduced with the locked package versions before this session was written; none is invented.

| Bug | Where | Effect | Used in |
| --- | --- | --- | --- |
| paho-mqtt 2.x callback arity | `I1820/mqtt/route.py` | The network thread dies on CONNACK; nothing is ever subscribed | Task 0 |
| Frozen default timestamp | `I1820/domain/log.py` | Every log from the dummy agent carries the import-time timestamp | Task 0 |
| Field lists shared across types | `I1820/things/base.py` | Every type reports every loaded type's settings and states | Task 2 |
| `last_seen` never refreshed | `I1820/discovery/discovery.py` | Agent time is stuck at the first ping | Task 3 |
| Reading an unset setting | `I1820/things/fields/setting.py` | KeyError becomes HTTP 500 | Task 3 |

The full interviewer guide, with traps, expected AI missteps and reference fixes, is [`docs/mia-cold-start.typ`](../docs/mia-cold-start.typ). Build it with `typst compile docs/mia-cold-start.typ`.

### Preparation

1. Fork the repository at `f0bd540` under a neutral name and remove `.github/` so nobody trusts the green CI badge. Keep the README; it documents the MQTT protocol and is part of the test.
2. Do not fix the paho bug. Task 0 depends on it.
3. Make the dummy agent's counter increase so "latest value" is observable. Keep the `I1820Log(...)` call without an explicit timestamp; that is the second trap.

```diff
--- dummy.py
-    i = 2
+    i = 0
     while True:
-        i = (i**2) % 1024
+        i += 1
         ping()
```

4. Pre-pull `mongo:5.0` and `emqx/emqx:5.0.0`, and prepare a Python 3.12 environment with the locked dependencies.
5. Verify the trap fires on the prep machine: run `main.py` and `dummy.py`, confirm `GET /agents` returns `{}` and that the paho `TypeError` traceback appears in the server terminal.

### Brief for the candidate

> You are joining a team that maintains MiA, a small IoT middleware. Agents publish device data over MQTT; users read and control devices over HTTP. You have never seen the code. You have two hours and an AI assistant of your choice.
>
> 1. Use the AI as much or as little as you like. We are interested in how you work, not in whether you use it.
> 2. Think aloud. Tell us what you are asking the AI and why, and what you believe after it answers.
> 3. After each task, commit and tag (`git tag task-N`). Unfinished is fine.
> 4. If something in a task is ambiguous, you may ask us. We may or may not answer.
> 5. Keep your AI transcript. We will collect it at the end.

### Task 0 — First light (15 min)

Start the stack (`docker compose up -d`, then `main.py` and `dummy.py`). Show us `GET /agents` listing an agent called `dummy` with one thing of type `dummy`, and `POST /things` for that thing returning its `chert` state, with a second call a few seconds later returning a higher value.

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap 1:** the handler is written for the paho-mqtt 1.x callback signature while the client is created with `CallbackAPIVersion.VERSION2`. On CONNACK paho calls `on_connect` with five arguments, the method accepts four, and the network thread dies with `TypeError: Handler.on_connect() takes 5 positional arguments but 6 were given`. The traceback is printed between rich log lines, then "Mia is up and running" appears and the HTTP server works normally. `GET /agents` returns `{}`.
- **Trap 2:** `I1820Log.__init__` has `timestamp=datetime.utcnow()` as a default argument, evaluated once at import. The dummy agent never passes a timestamp, so every log carries the same one and "latest" stops meaning anything.
- **Expected AI missteps:** blaming Docker networking, the tenant ID or the broker for the empty agent list; proposing to sort by `_id` instead of `time` for the stuck value, which hides the symptom and works.
- **Watch for:** whether they read the server terminal or paste "agents is empty" without the traceback; whether they test a theory before applying it; whether they ask *why* the timestamps are equal or take the working patch; whether they re-run the full loop after each fix.

</details>

### Task 1 — Trace a log (15 min)

In at most ten lines, write down what happens between an agent publishing to `I1820/main/log/send` and a user reading that value through `POST /things`. Name the file and function at each step. Then answer: where does the value live between those two moments, and what happens if the agent never sent a ping?

<details>
<summary><b>Interviewer notes</b></summary>

- **No trap.** This is the baseline for how the candidate uses the AI for comprehension, and it primes them for Task 2.
- **Correct chain:** `Handler.on_log` validates against `log.json` via `I1820Log.from_json`; `Things.get(type)` imports `I1820.things.models.<type>` on demand, the metaclass having registered the class at import; `get_thing` looks the instance up in `registered_things`, which `DiscoveryService.ping` populated, so without a ping the log is dropped with a warning; `setattr` hits `State.__set__`, which goes through the `LogController` singleton (wired in `main.py`) to `MongodbLogAppender.create`; on read `State.__get__` calls `retrieve_last`. The value lives only in Mongo; the thing instance holds nothing.
- **Expected AI missteps:** claiming state is cached on the instance; describing `Things.get` as reading a registry file; missing that `LogController` is a singleton wired from `main.py`.
- **Watch for:** whether they ask for line references and open the files, and whether they use the runtime to confirm, for example by stopping `dummy.py` and restarting the server to see the log-drop warning.

</details>

### Task 2 — Dimmable lamp (30 min)

Add a thing type `dimmable_lamp`. It has everything a `lamp` has, plus a `brightness` setting (integer, 0 to 100) and a `power` statistic. `GET /models/dimmable_lamp` must describe exactly its fields: settings `on` and `brightness`, statistics `power`, nothing else. `PUT /things` with `brightness: 60` must publish a notification on `I1820/main/configuration/request`; show it with any MQTT subscriber. Tests cover both.

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** `class DimmableLamp(Lamp)` with two new fields is the correct shape. But `AbstractThing.__new__` appends fields to the lists declared on `Thing` itself, so every type shares them. Once `lamp`, `curtain` and `temperature` are imported, the endpoint reports settings `on, height, brightness` and states `temperature`. The same bug makes `POST /things` with an empty `states` list return 500 on any actuator. The proper fix is in the metaclass: give each class its own lists, seeded from its bases. A filter in `to_dict` hides the symptom for one endpoint only.
- **Smaller trap:** the file must be `models/dimmable_lamp.py`, because `Things.get` imports by module name and then looks up the snake-cased class name. `dimmablelamp.py` imports fine and then raises a misleading `ThingTypeNotImplementedException`.
- **Expected AI missteps:** writing the class, a test that checks `issubclass`, and stopping; putting brightness validation in the HTTP handler although `Setting.__set__` publishes before anything can validate; inventing a mock for the MQTT service although `test_thing.py` already shows the `DummyMQTTService` pattern.
- **Watch for:** whether they call `GET /models/dimmable_lamp` after the AI says it is done; whether they widen the investigation to existing types; root cause versus symptom; whether they notice the duplicated mock.

</details>

### Task 3 — Two bug reports (25 min)

Two tickets from operations. Fix both, with regression tests that fail before the fix and pass after.

- **OPS-118.** "The `time` field in `GET /agents` never changes for an agent that has been pinging for hours. Our stale-agent alert never fires."
- **OPS-121.** "Reading a lamp through `POST /things` with `settings: ["on"]` returns HTTP 500 if nobody has switched it yet. We expect a clean answer."

<details>
<summary><b>Interviewer notes</b></summary>

- **Trap:** in `DiscoveryService.ping` the re-ping branch does `agent.last_seen = now` on the incoming object, which is then discarded. Two lines above sits `if self._agents[agent.ident] == agent.things`, comparing an `Agent` to a `set`, always false. The AI frequently fixes that comparison and reports success while `last_seen` is still stale. The existing `test_discovery.py` passes throughout.
- **Decision:** `Setting.__get__` does an unguarded dict lookup. The README states that settings must be set before they are read, so the contract is a real choice: return `null` like `Statistic` does, or 404, or 400. The candidate should notice this is not theirs to decide alone. Returning `None` also needs the return annotation changed to `T | None`, or mypy fails.
- **Expected AI missteps:** a regression test that constructs two pings without advancing time and therefore passes in both states; wrapping the HTTP handler in a try/except that hides the next KeyError too.
- **Watch for:** whether they run the new test before the fix; how they catch the wrong-line fix, by diff, by running the server, or not at all; whether they surface the contract question.

</details>

### Task 4 — Many devices (stretch)

`POST /things` accepts a list of `device_id`s but returns a single flat object, so only the last device's values survive. Make it return per-device results without breaking the single-device call.

<details>
<summary><b>Interviewer notes</b></summary>

- **No trap.** A design task with no correct answer to verify against. Watch whether they consider backward compatibility and ask about the response shape, or let the AI pick one.

</details>

### Debrief (15 min)

1. Draw the path of one log message from the agent to Mongo. Which file registers the thing type, and when?
2. Which of the AI's claims turned out to be wrong today? How did you find out?
3. Was there a moment you stopped trusting the AI? What triggered it?
4. Which change would you be least comfortable merging, and why?
5. Two MiA servers with different tenants share one broker. What still leaks between them?
6. What did you wish the AI could see or do that it could not?

### Collect

Git history with the `task-N` tags including the final working tree, the AI transcript, the screen recording with the think-aloud audio, and the interviewer sheet with timestamped notes.

### Caveats

- The five bugs were reproduced in a Python 3.12 environment with the locked package versions (paho-mqtt 2.1.0, sanic 24.6.0, pymongo 4.10.1). The MQTT failure was reproduced against an in-process broker, not EMQX.
- The exact value Mongo returns among equal timestamps was not verified. The trap holds because the "higher value on the second call" criterion cannot be reliably met, but confirm the observed behaviour on the prep machine.
- The time boxes assume a prepared environment. If the candidate must install images or dependencies, add 20 minutes to Task 0 and treat that time as noise.
