---
name: git-commit
description: Lab git commit conventions
---

Format:

```
area: description

body

[Tested locally:

  % command
  result]

Issue
```

Rules:
- lines are up to 80 characters long
- subject is imperative, lowercase, without trailing period
- body is concise: why in one or two sentences, then what changed
- prove local testing under `Tested locally:`: the actual command, kept
  simple, and its short result, indented by 2 spaces with a `% ` prompt
- reference the issue or NO_ISSUE, close issue with Fix

Process:

1. Stage files explicitly. Do not use `git add -A`.
2. Show the draft message for review before committing.
3. If a hook fails, fix the cause and commit again. Never use `--no-verify`.
4. Do not push unless asked.
