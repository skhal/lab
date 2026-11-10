---
name: git-commit
description: Lab git commit conventions
---

Format:

```
area: description

body

[Issue #123|Fix #123|NO_ISSUE]
```

Rules:
- lines are up to 80 characters long
- subject is imperative, lowercase, without trailing period
- body has a short description of why and what the change is
- include a proof of the fix if available
- reference the issue or NO_ISSUE

Process:

1. Stage files explicitly. Do not use `git add -A`.
2. Show the draft message for review before committing.
3. If a hook fails, fix the cause and commit again. Never use `--no-verify`.
4. Do not push unless asked.
