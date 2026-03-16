# AI Collaboration Workflow

## Purpose

This directory contains reusable prompts, templates, and transition artifacts
for low-latency AI-assisted development.

## Structure

- `prompts/` — static prompts reused across context switches
- `templates/` — reusable file templates
- `sessions/` — generated handoff artifacts for the next chat
- `docs/` — durable project knowledge such as `MasterPlan.md`
- `user/` — user-specific preferences used to keep interaction consistent

## Workflow

1. Work in a chat until context size, latency, or topic drift becomes a problem.
2. Use `prompts/compression_prompt.txt` for mid-chat cleanup when helpful.
3. Use `prompts/transition_prompt.txt` when preparing to switch to a new chat.
4. Update `docs/MasterPlan.md` as the durable project record.
5. Generate fresh files in `sessions/`:
   - `checkpoint.txt`
   - `bootstrap_prompt.txt`
   - `current_task.txt`
   - optionally `open_questions.txt`
6. Start a new chat and paste `sessions/bootstrap_prompt.txt`.
7. Then paste only the current task, relevant files, and minimal failure output.

## Guidelines

- Keep durable project truth in `docs/`, not in chat-only artifacts.
- Keep handoff artifacts short and factual.
- Treat newly pasted code, tests, and errors as the current source of truth.
- Prefer compact state capsules over long historical summaries.