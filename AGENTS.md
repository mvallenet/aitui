# AGENTS.md

## Project

`aitui` is a reusable Go library for CLI developers who want to plug AI
backends into their tools through one stable interface.

The library should support both API-based providers and official CLI harnesses
such as Codex CLI. Do not reimplement subscription authentication, scrape
private credentials, or store raw API keys in library-managed state.

## Repository Setup

- Use the Go version declared in `go.mod` once it exists.
- Install dependencies with:
  - `go mod download`
- Keep generated files out of source control unless they are required for
  consumers of the library.

## Common Commands

- Format code:
  - `gofmt -w .`
- Check packages:
  - `go test ./...`
- Run focused tests:
  - `go test ./path/to/package -run TestName`
- Vet when changing public APIs, concurrency, or command execution:
  - `go vet ./...`

## Code Style

- Write idiomatic Go with small packages and clear ownership boundaries.
- Keep public APIs simple, stable, and easy to mock in downstream CLIs.
- Prefer deterministic, testable code over hidden global state.
- Return clear, human-readable errors that include useful context without
  leaking secrets.
- Use interfaces at package boundaries when they simplify testing or provider
  substitution; avoid abstractions that only prepare for imagined features.
- Keep provider-specific behavior behind narrow adapters.
- Do not over-engineer the initial version.

## Provider Guidance

- API providers should accept credentials through caller-owned configuration,
  environment variables, or explicit transport setup.
- CLI harness providers should call official tools and rely on their supported
  authentication flows.
- Never persist raw API keys, session tokens, or provider secrets from inside
  this library.
- Make command execution predictable: explicit args, bounded contexts, captured
  stderr, and errors that explain what failed.

## Testing Expectations

- Add table-driven unit tests for parsing, request construction, adapter
  behavior, and error cases.
- Avoid live provider calls in default tests.
- Gate integration tests behind explicit environment variables or build tags.
- Prefer fake transports and fake command runners for provider tests.

## Future Docs

Use this file as a table of contents. Put durable explanations in:

- `docs/architecture.md` for package layout, core interfaces, and provider
  adapter design.
- `docs/testing.md` for unit, integration, and fixture strategy.
- `docs/debugging.md` for provider troubleshooting, command logs, and safe
  redaction practices.

## Agent Notes

- Read existing package docs and tests before changing public behavior.
- Keep changes scoped to the requested feature or bug.
- Update docs when behavior, setup, or public APIs change.
- Preserve user-owned work in the repository; do not revert unrelated changes.
