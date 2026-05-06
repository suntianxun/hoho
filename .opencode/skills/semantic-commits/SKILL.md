---
name: semantic-commits
description: Use when finishing code changes, making a git commit, or being asked to save progress to ensure semantic commit prefixes.
---

# Semantic Commits

## Overview
Ensure all git commits follow semantic commit guidelines (Conventional Commits) with proper prefixes like feat, fix, docs, style, refactor, perf, test, build, ci, chore, or revert.

## When to Use
- Whenever you finish code changes and are about to commit.
- Whenever the user explicitly asks you to commit code.

## Core Pattern
Always prefix the commit message with the appropriate type.

```bash
# Good
git commit -m "feat: add user authentication"
git commit -m "fix: resolve memory leak in worker"

# Bad
git commit -m "added user auth"
git commit -m "fixed bug"
```

## Quick Reference
| Type | Use Case |
|---|---|
| feat | A new feature |
| fix | A bug fix |
| docs | Documentation only changes |
| style | Changes that do not affect the meaning of the code (white-space, formatting, etc) |
| refactor | A code change that neither fixes a bug nor adds a feature |
| perf | A code change that improves performance |
| test | Adding missing tests or correcting existing tests |
| build | Changes that affect the build system or external dependencies |
| ci | Changes to our CI configuration files and scripts |
| chore | Other changes that don't modify src or test files |
| revert | Reverts a previous commit |

## Implementation
When committing, analyze the changes staged. Choose the most appropriate prefix. If there are multiple types of changes, try to separate them into different commits if logically distinct, or pick the most significant type.

## Common Mistakes
- Forgetting the colon after the prefix.
- Capitalizing the prefix (e.g., `Feat:` instead of `feat:`).
- Not using imperative mood in the subject line (e.g., `fix: fixed bug` instead of `fix: resolve bug`).