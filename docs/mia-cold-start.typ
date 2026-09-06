#set document(
  title: "AI Session: MiA Cold Start",
  author: "Parham Alvani",
)

// ---------------------------------------------------------------- palette

#let ink = rgb("#161310")        // near-black
#let accent = rgb("#c2410c")     // orange
#let accent-tint = rgb("#fdf0e7")
#let warn = rgb("#9b2c2c")       // a darker red, so it is not mistaken for the orange
#let warn-tint = rgb("#fdf2f2")
#let grey = rgb("#6e6762")
#let rule = rgb("#e3ddd8")

// ---------------------------------------------------------------- page

#set page(
  paper: "a4",
  margin: (top: 1.7cm, bottom: 1.6cm, left: 1.8cm, right: 1.8cm),
  footer: context {
    set text(size: 8pt, fill: grey)
    grid(
      columns: (1fr, auto),
      align(start)[#link("https://github.com/1995parham-teaching/interviews")[interviews] repository],
      align(end)[#counter(page).display("1 of 1", both: true)],
    )
  },
)

#set text(font: "Charter", size: 10pt, fill: ink, lang: "en", hyphenate: false)
#set par(justify: true, leading: 0.68em, spacing: 0.85em)

#set enum(numbering: "1.")
#set heading(numbering: none)
#show heading.where(level: 1): it => block(above: 1.3em, below: 0.5em)[
  #set text(font: "Helvetica Neue", size: 12pt, weight: 700, fill: accent)
  #it.body
]
#show heading.where(level: 2): it => block(above: 1.1em, below: 0.4em)[
  #set text(font: "Helvetica Neue", size: 10pt, weight: 700)
  #it.body
]

#set table(
  stroke: (x, y) => if y == 0 { (bottom: 0.8pt + accent.lighten(40%)) } else { (bottom: 0.4pt + rule) },
  inset: (x: 6pt, y: 5.5pt),
  fill: (x, y) => if y == 0 { accent-tint } else { none },
)
#show table.cell.where(y: 0): set text(font: "Helvetica Neue", weight: 700, fill: accent, size: 8.5pt)
#show table: set par(justify: false)
#show table: set text(size: 9pt)

#show raw: set text(font: "Menlo", size: 8.5pt)
#show raw.where(block: true): it => block(
  fill: accent-tint.lighten(50%),
  stroke: 0.4pt + rule,
  radius: 2pt,
  inset: (x: 8pt, y: 6pt),
  width: 100%,
  above: 0.7em,
  below: 0.8em,
)[#it]
#show link: it => text(fill: accent)[#underline(it, offset: 1.5pt, stroke: 0.4pt)]

#let callout(color, bg, title, body) = block(
  fill: bg,
  stroke: (left: (paint: color, thickness: 2.5pt)),
  radius: (right: 3pt),
  inset: (x: 9pt, y: 8pt),
  width: 100%,
  above: 0.9em,
  below: 0.9em,
  breakable: true,
)[
  #set par(justify: false)
  #text(font: "Helvetica Neue", weight: 700, size: 9pt, fill: color)[#title] #h(0.35em) #body
]
#let note(t, b) = callout(accent, accent-tint, t, b)
#let trap(t, b) = callout(warn, warn-tint, t, b)

// the statement handed to the candidate
#let brief(body) = block(
  stroke: 0.6pt + rule,
  radius: 3pt,
  inset: (x: 10pt, y: 8pt),
  width: 100%,
  above: 0.8em,
  below: 0.9em,
  breakable: false,
)[
  #set par(justify: false)
  #text(font: "Helvetica Neue", weight: 700, size: 8pt, fill: grey, tracking: 0.06em)[BRIEF FOR THE CANDIDATE] #v(0.3em) #body
]

// task heading with its time budget
#let task(n, title, time) = block(above: 1.6em, below: 0.4em)[
  #grid(
    columns: (auto, 1fr, auto),
    column-gutter: 8pt,
    align: (start, start, end),
    text(size: 8.5pt, fill: accent, font: "Menlo")[TASK #n],
    text(font: "Helvetica Neue", size: 11pt, weight: 700)[#title],
    text(font: "Helvetica Neue", size: 8.5pt, fill: grey)[#time],
  )
  #line(length: 100%, stroke: 0.5pt + rule)
]

// checklist with tick boxes
#let tick = box(width: 7.5pt, height: 7.5pt, radius: 1.5pt, stroke: 0.7pt + accent, baseline: 1pt)
#let checklist(..items) = block(above: 0.8em, below: 0.8em, width: 100%)[
  #grid(
    columns: (8pt, 1fr),
    column-gutter: 7pt,
    row-gutter: 0.62em,
    ..items.pos().map(i => (tick, i)).flatten()
  )
]

// ---------------------------------------------------------------- document header

#block(below: 1.1em)[
  #block(below: 4pt)[
    #text(font: "Helvetica Neue", size: 9pt, weight: 700, fill: accent, tracking: 0.04em)[INTERVIEWER GUIDE]
  ]
  #block(below: 11pt)[
    #text(font: "Helvetica Neue", size: 17pt, weight: 700)[AI Session: MiA Cold Start]
  ]
  #line(length: 100%, stroke: 1.6pt + accent)
]

= What this session measures

A two-hour session on a codebase the candidate has never seen, with one twist: we expect them to use an AI assistant. Whether they finish does not matter. What matters is how they work with the assistant inside unfamiliar code, and whether they notice when the assistant is wrong.

= Why this codebase

#link("https://github.com/I1820/mia")[MiA] is a small IoT middleware. Agents publish device data over MQTT; users read and control devices over HTTP. The whole project is about 1,900 lines of Python across 71 files, with six tests that all pass.

It is small enough to read in an hour, but it is built on a few Python features most of us only half-remember: a metaclass that registers every `Thing` subclass by name; descriptors with side effects (assigning to a `Setting` publishes an MQTT message, assigning to a `State` writes to Mongo); singleton controllers whose dependencies are set as class attributes in `main.py`; and on-demand loading of device types with `importlib` by module name. An AI assistant explains all of this fluently and gets the details wrong. That is exactly what we want.

The project also has a real runtime: an MQTT broker, a Mongo database, a dummy agent and an HTTP API. So the candidate can actually try what the assistant claims instead of only reading the code. Whether they do is the main thing we watch for.

The `main` branch at commit `f0bd540` carries five latent bugs that CI does not catch. All of them were reproduced with the locked library versions before this session was written; none is invented.

#table(
  columns: (1.1fr, auto, 1.4fr, auto),
  align: (left, left, left, center),
  table.header([Bug], [Where], [Effect], [Task]),
  [paho-mqtt 2.x callback arity], [`I1820/mqtt/route.py`], [The network thread dies on CONNACK; no topic is ever subscribed], [0],
  [Frozen default timestamp], [`I1820/domain/log.py`], [Every log from the dummy agent carries the same timestamp: the moment of import], [0],
  [Field lists shared across types], [`I1820/things/base.py`], [Every type reports the fields of every loaded type], [2],
  [`last_seen` never refreshed], [`I1820/discovery/discovery.py`], [The agent's time stays at the first ping], [3],
  [Reading an unset setting], [`I1820/things/fields/setting.py`], [`KeyError` becomes HTTP 500], [3],
)

= Format

/ Individual, 120 minutes: four timed tasks plus a stretch task, then a 15 minute debrief.
/ Their own AI: any tool is allowed; Claude Code, Cursor, Copilot, or a plain chat window. Note which one it was and whether it could run commands by itself. A chat window forces the candidate to be the hands; with an agentic tool the question becomes whether the candidate reads what the tool did.
/ Observed: one interviewer with a notes sheet, plus screen recording. Without an observer, the AI transcript and the git history keep most of what you need.
/ Commit after every task: the candidate tags the end of each task (`task-0`, `task-1`, ...). Unfinished work is fine and is evidence in itself.
/ Tasks are handed out one at a time: later tasks refer to bugs the earlier tasks may have exposed. Handing them out early spoils the result.
/ The environment is prepared: images pulled, dependencies installed. Time spent on setup tells us nothing.

= Preparation

Create a fork under a neutral name (for example `mia-challenge`). The whole preparation takes under an hour.

#checklist(
  [Fork at commit `f0bd540` and delete the `.github/` directory so nobody trusts the green CI badge. Keep the README; it documents the MQTT protocol and is part of the test.],
  [Do *not* fix the paho bug. Task 0 depends on it.],
  [Make the dummy agent's counter increase so that "the latest value" means something. Keep the `I1820Log(...)` call without a timestamp; that is the second trap.],
  [Pull the `mongo:5.0` and `emqx/emqx:5.0.0` images and prepare a Python 3.12 environment with the locked dependencies.],
  [Confirm the trap fires on the prep machine: run `main.py` and `dummy.py`. `GET /agents` must return `{}` and the paho `TypeError` traceback must show up in the server terminal.],
  [Print this guide and write each task's brief on a separate card.],
)

The change to `dummy.py`:

```diff
-    i = 2
+    i = 0
     while True:
-        i = (i**2) % 1024
+        i += 1
         ping()
```

= What you hand to the candidate

#note[Read verbatim:][
  You are joining a team that maintains MiA, a small IoT middleware. Agents publish device data over MQTT; users read and control devices over HTTP. You have never seen this code. You have two hours and an AI assistant of your choice.

  1. Use the AI as much or as little as you like. We care about how you work, not about whether you use it.
  2. Think aloud. Tell us what you ask the AI and why, and what you conclude after it answers.
  3. After each task, commit and tag (`git tag task-N`). Unfinished is fine.
  4. If something in a task is unclear, you may ask us. We may or may not answer.
  5. Keep your AI transcript. We collect it at the end.
]

#task("0", "First light", "15 minutes")

#brief[
  Bring the stack up (`docker compose up -d`, then `main.py` and `dummy.py`). Show us that:

  - `GET /agents` lists an agent called `dummy` with one device of type `dummy`.
  - `POST /things` for that device returns its `chert` state, and a second call a few seconds later returns a *higher* value.
]

#trap[Trap 1 · The server says it is up, but it hears nothing.][
  The handler is written for the paho 1.x callback signature, but the client is created with `CallbackAPIVersion.VERSION2`. On CONNACK, paho calls `on_connect` with five arguments, the method accepts only four, and the network thread dies with:

  ```
  TypeError: Handler.on_connect() takes 5 positional arguments but 6 were given
  ```

  The traceback is printed between the rich log lines, then "Mia is up and running" appears and the HTTP server works normally. `GET /agents` answers `{}`. No topic is ever subscribed.
]

#trap[Trap 2 · The value never moves.][
  The `I1820Log` constructor has the default argument `timestamp=datetime.utcnow()`, which is evaluated once, at import time. The dummy agent never sends a timestamp, so all of its logs carry the same one. Mongo's `find_one(sort=time desc)` picks one of the ties, and "latest" no longer means anything. The second criterion cannot be met until the candidate understands this.
]

#note[What to watch for:][
  - Do they read the server terminal, or just write "agents is empty" and leave the traceback out of the prompt?
  - When the AI blames Docker networking, the tenant, or the broker, do they test the theory first or apply the suggested change as is?
  - For trap 2, the AI usually proposes sorting by `_id` instead of `time`. That only hides the symptom, but it works. Does the candidate ask *why* the timestamps are equal, or accept the patch?
  - After each fix, do they run the whole path again, or declare success at the first non-empty answer?
]

#task("1", "The path of a log", "15 minutes")

#brief[
  In at most ten lines, write down what happens from the moment an agent publishes on `I1820/main/log/send` to the moment a user reads that value through `POST /things`. Name the file and the function at each step. Then answer two questions: where is the value kept between those two moments? And what would happen if the agent had never sent a ping?
]

This task has no trap. Its purpose is to see how the candidate uses the AI to understand code, and it also prepares them for Task 2, which needs an understanding of the metaclass. The correct path:

+ `Handler.on_log` in `mqtt/route.py` validates the input with `I1820Log.from_json` against `log.json`.
+ `Things.get(type)` in `things/base.py` imports the module `I1820.things.models.<type>` on the spot, and the metaclass registers the class right there.
+ `get_thing(agent, device)` finds the instance in `registered_things`; that dict was filled by `DiscoveryService.ping`. Without a ping, `ThingNotFoundException` is raised and the log is dropped with a warning.
+ `setattr(thing, name, {...})` reaches `State.__set__`. That calls the `LogController` singleton, whose `log_service` is set as a class attribute in `main.py`, and ends in `MongodbLogAppender.create`.
+ On read, `getattr` reaches `State.__get__`, which reads from Mongo with `retrieve_last`. The value exists only in Mongo; the thing instance holds nothing.

#note[What to watch for:][
  - Do they accept the AI's explanation as is, or ask it for line numbers and open the files themselves?
  - Common AI mistakes on this task: it says state is cached on the instance; it describes `Things.get` as reading a registry file; it misses that `LogController` is a singleton configured in `main.py`. Count how many the candidate catches.
  - Do they use the runtime to make sure? For example, stop `dummy.py` and restart the server to see the warning about the dropped log.
]

#task("2", "Dimmable lamp", "30 minutes")

#brief[
  Add a device type `dimmable_lamp`. It has everything a `lamp` has, plus a `brightness` setting (integer, 0 to 100) and a `power` statistic. Acceptance criteria:

  - `GET /models/dimmable_lamp` returns exactly its own fields: settings `on` and `brightness`, statistic `power`, and nothing else.
  - `PUT /things` with `brightness: 60` publishes a notification on `I1820/main/configuration/request`. Show it with any MQTT subscriber.
  - Write tests for both.
]

#trap[Trap · The subclass is right, but the endpoint still answers wrong.][
  `class DimmableLamp(Lamp)` with two new fields is exactly what it should be. But `AbstractThing.__new__` collects fields with `getattr(instance, field_name).append(value)`, and the lists it appends to are the ones defined on `Thing` itself. In other words, all types share one list. Once `lamp`, `curtain` and `temperature` have been imported, the endpoint returns:

  ```json
  {"type": "dimmable_lamp", "settings": ["on", "height", "brightness"],
   "states": ["temperature"], "statistics": ["power"], "events": []}
  ```

  The same bug breaks `POST /things` with an empty `states` list on any actuator, because the handler then calls `getattr(lamp, "temperature")` and returns 500. The right fix is in the metaclass: each class gets its own lists, seeded with the fields of its base classes. Filtering in `to_dict` only hides the symptom for one endpoint.

  A smaller trap: the file must be `models/dimmable_lamp.py`, because `Things.get` imports by module name and then looks the class up by its snake_case name. `dimmablelamp.py` imports without error and then raises a misleading `ThingTypeNotImplementedException`.
]

#note[What to watch for:][
  - When the AI says it is done, do they actually call `GET /models/dimmable_lamp`? In most runs the AI writes the class, maybe a test for `issubclass`, and stops there.
  - When they see the extra fields, do they check the existing types too, or fix only the new one?
  - Root cause or symptom? A metaclass fix, a filter in `to_dict`, or a hardcoded list?
  - Where does the brightness range validation go? The AI usually puts it in the HTTP handler. Does the candidate ask whether the descriptor is a better place, given that `Setting.__set__` publishes before any validation can happen?
  - For the notification test, `test_thing.py` already shows the `DummyMQTTService` pattern. Does the AI find it or build a fresh mock? Does the candidate notice the duplication?
]

#task("3", "Two bug reports", "25 minutes")

#brief[
  Two tickets from the operations team. Fix both, and for each write a regression test that fails before the fix and passes after it.

  - *OPS-118*: "The `time` field in `GET /agents` never changes for an agent that has been pinging for hours. Because of that, our stale-agent alert never fires."
  - *OPS-121*: "Reading a lamp with `POST /things` and `settings: ["on"]` returns HTTP 500 if nobody has switched it on or off yet. We want a clean answer."
]

#trap[Trap · Two wrong lines next to the right one.][
  In `DiscoveryService.ping`, on the branch for repeated pings, `agent.last_seen = now` is done on the incoming object, which is thrown away right after. The stored agent keeps its first timestamp. Two lines above there is also `if self._agents[agent.ident] == agent.things`, which compares an `Agent` with a `set` and is always false. The AI often "fixes" that comparison and reports success while `last_seen` is still stale. The existing `test_discovery.py` passes the whole time.

  In OPS-121, `Setting.__get__` reads from the dict with no guard. The README says settings must be set before they are read, so there is a real decision about the API contract here: return `null` like `Statistic` does, or 404, or 400. The candidate should realise this decision is not theirs alone. If they return `None`, the return type must also become `T | None`, or mypy fails in CI.
]

#note[What to watch for:][
  - Do they run the new test *before* the fix to see it fail? A test the AI writes with two pings in a row and no time advance passes both before and after the fix.
  - Do they notice that the AI changed the wrong line? How: by reading the diff, by running the server, or not at all?
  - Do they raise the API contract question with the interviewer, or let the AI decide quietly?
  - Is the OPS-121 fix made in the descriptor, or does the AI wrap the handler in a try/except that also hides the next `KeyError`?
]

#task("4", "Several devices (stretch)", "if time remains")

#brief[
  `POST /things` accepts a list of `device_id`s but returns one flat object, so only the last device's values survive in it. Make the result come back per device, without breaking the single-device call.
]

This task is more design than bug hunting, and it has no trap. We want to see whether the candidate considers compatibility with the previous version and asks about the response shape, or lets the AI pick one. It fills the remaining time and shows how they behave when there is no "right" answer to check against.

= Scoring: the same four competencies

Every task is scored with the *same four competencies* and the *same scale* as the evaluation document: *Weak · Borderline · Strong · Not assessed*. The AI does not change the competencies; it only changes what counts as evidence for each.

#table(
  columns: (auto, 1fr, 1fr),
  align: (left, left, left),
  table.header([Competency], [Weak], [Strong]),
  [Data structures, algorithms and problem solving],
  [Patches where the symptom shows up (`_id`, a filter in `to_dict`); at the debrief cannot say where state is kept or what the metaclass does],
  [Finds the root cause of each bug and checks its effect elsewhere; draws the data path unaided],
  [Coding],
  [Lets the AI rewrite the metaclass, add pydantic, or migrate to FastAPI; their test passes before the fix too],
  [Keeps the change limited to the task; writes a failing test first; follows the patterns already in the repository],
  [Effective communication],
  [Prompts like "agents is empty, why"; accepts the AI's explanation as is; the API contract is decided quietly by the AI],
  [Prompts carry the traceback, the file name and the expectation; rejects a wrong explanation with evidence; states a design decision explicitly and asks, or names the assumption],
  [Speed and getting to a result],
  [Tries the AI's suggestions one by one until one works; never calls the endpoint or runs the tests],
  [Reproduces the problem before the fix and checks the result after it, reads the diff, and still reaches working code within the time budget],
)

For each task, also record three simple counts: *files the candidate opened*, *files the AI opened*, and *AI claims accepted without checking*. These three numbers separate "reads code" from "reads about code" better than any score.

= Collecting data, and the debrief

/ What to collect: the git history with the `task-N` tags and the final working tree; the AI transcript (a Claude Code session export, the Cursor chat, or the text of the chat window); the screen recording with the candidate's audio; the interviewer's sheet with timestamped notes.

/ Debrief, 15 minutes:
  + Draw the path of one log message from the agent to Mongo. Which file registers the device type, and when?
  + Which of the AI's claims turned out to be wrong today? How did you find out?
  + Was there a point where you lost trust in the AI? What caused it?
  + Which change would you merge with the least confidence, and why?
  + Two MiA servers with different tenants share one broker. What do they still have in common? (Answer: nothing over MQTT; the type registry and the controller singletons are per process too. A candidate who thinks it through and says so has the right model of the system.)
  + What did you wish the AI could see or do that it could not?

= Reference fixes

A short version for the interviewer. All five fixes were applied together to a copy of the repository: the six existing tests and mypy pass, and every bug in the table returns to the correct behaviour.

== Task 0 · the paho callback

```python
# I1820/mqtt/route.py
def on_connect(self, client, userdata, flags, reason_code, properties):
    ...
```

== Task 0 · the frozen timestamp

```python
# I1820/domain/log.py
def __init__(self, type, device, states, agent, timestamp=None):
    ...
    self.timestamp = timestamp or datetime.datetime.now(datetime.UTC)
```

== Task 2 · a separate field list per class

```python
# I1820/things/base.py, inside AbstractThing.__new__ after the abstract check
for attr in ("statistics", "settings", "events", "states"):
    inherited = [f for base in bases for f in getattr(base, attr, [])]
    setattr(instance, attr, list(inherited))
for _, value in namespace.items():
    if isinstance(value, Field):
        getattr(instance, value.field_name).append(value)
```

== Task 3 · last_seen

```python
# I1820/discovery/discovery.py, in the else branch of ping()
stored = self._agents[agent.ident]
stored.last_seen = datetime.datetime.now(datetime.UTC)
if stored.things == agent.things:
    return
```

== Task 3 · the unset setting

```python
# I1820/things/fields/setting.py
def __get__(self, obj: Thing, _) -> T | None:
    return self.storage.get((obj.agent_id, obj.device_id))
```

Returning `None` is consistent with `Statistic` and with the shape of the read handler's response. A 404 is defensible too, if the team prefers it. What matters is that the behaviour was chosen, not inherited from a `KeyError`.

= Caveats

- The five bugs were reproduced in a Python 3.12 environment with the locked versions (paho-mqtt 2.1.0, sanic 24.6.0, pymongo 4.10.1). The MQTT failure was reproduced against an in-process broker, not EMQX.
- Which value Mongo returns among equal timestamps was not checked. The trap still holds, because the "higher value on the second call" criterion cannot be met reliably. Even so, look at the real behaviour on the prep machine and adjust the brief if needed.
- The timing assumes a prepared environment. If the candidate has to install images or dependencies, add 20 minutes to Task 0 and leave that time out of the evaluation.
