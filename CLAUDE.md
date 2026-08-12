# Conventions

## Colours

**Orange and black are the primary colours.** Use them for anything produced for this repository —
documents, diagrams, slides, generated pages. Do not introduce blue as an accent.

The palette in `docs/designing-interview-questions.typ` is the reference:

| Role | Value | Notes |
| --- | --- | --- |
| Ink | `#161310` | Near-black, slightly warm. Body text. |
| Accent | `#c2410c` | Orange. Headings, rules, table headers, links. |
| Accent tint | `#fdf0e7` | Table header fill, note callouts. |
| Warning | `#9b2c2c` | Deliberately a darker red so it is not mistaken for the orange. |
| Warning tint | `#fdf2f2` | |
| Grey | `#6e6762` | Secondary text, footers. Neutral, not blue-grey. |
| Rule | `#e3ddd8` | Hairlines and separators. Neutral, not blue-grey. |

## Persian documents

- Font: Vazirmatn. Set `lang: "fa"` and `dir: rtl`.
- Code, identifiers and numbers stay ASCII and LTR; only comments beside them are Persian.
- Use نیم‌فاصله (`می‌شود`, `داده‌ساختار`).
- Persian digits for page numbers and heading numbers (`counter(page).display("۱ از ۱", both: true)`
  renders both the current page and the total through the Persian pattern — a bare integer does not).
- No cover page or table of contents on a document of only a few pages.
- No horizontal rule under section headings. Colour and weight carry the hierarchy.

## Levels

Questions and problems are tagged **1 Simple / 2 Medium / 3 Complex**. A level describes the nature
of the work, never a specialty or job title. Level 1 is not the screening round and Level 3 is not a
system design assessment.

## Scoring

Score against the four competencies in the evaluation document (see
`skills/coding-interview-report.skill`), on its scale — ضعیف / مرزی / قوی / سنجیده نشد. Never invent
a per-question scale.

## Build

```
typst compile docs/designing-interview-questions.typ
```

The rendered PDF is committed; `.gitignore` excludes `*.pdf` everywhere except `docs/`. Rebuild and
commit both files together so they cannot drift.

## Commits

Do not add a `Co-Authored-By` trailer.
