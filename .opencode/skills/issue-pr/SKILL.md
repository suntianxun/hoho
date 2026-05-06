---
name: issue-pr
description: Use when asked to issue a PR, create a pull request, or submit code for review to ensure an issue is created first and linked.
---

# Issue PR Workflow

## Overview
When asked to create a pull request, always create a GitHub issue first describing the changes, then create the pull request and link it to the issue.

## When to Use
- Whenever the user says "issue PR", "create PR", "make a pull request".

## Core Pattern
1. Create the issue using the `gh` CLI.
2. Extract the issue number from the output.
3. Create the pull request using the `gh` CLI, ensuring you reference the issue number in the PR body (e.g., `Resolves #123` or `Fixes #123`).

## Implementation
Use the bash tool to run GitHub CLI commands.

```bash
# 1. Create issue
gh issue create --title "Describe the issue" --body "Detailed description of what is being fixed or added."

# Note the issue number from the output, e.g., https://github.com/org/repo/issues/42 -> Issue #42

# 2. Create branch, commit, push if not already done.
# Make sure commits follow semantic-commits!

# 3. Create PR linking the issue
gh pr create --title "feat: descriptive title" --body "Resolves #42. Description of changes."
```

## Common Mistakes
- Creating a PR without creating an issue first.
- Forgetting to link the issue in the PR body using a closing keyword like `Resolves #NUM` or `Fixes #NUM`.