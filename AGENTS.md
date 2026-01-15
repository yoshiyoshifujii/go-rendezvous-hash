# Repository Guidelines

## Project Structure & Module Organization
- Core implementation lives at repo root: `rendezvous.go`, `node.go`, `hasher.go`, `types.go`.
- Tests are colocated: `rendezvous_test.go`.
- Examples should live under `examples/` if added in the future.

## Build, Test, and Development Commands
- `go test ./...` runs all tests.
- `go vet ./...` performs basic static analysis.
- `gofmt -w *.go` formats the Go sources.

## Coding Style & Naming Conventions
- Use standard Go formatting (`gofmt`) and keep indentation at tabs as enforced by Go tooling.
- Package name is `rendezvoushash`; keep new packages short and lower-case if added later.
- Exported identifiers should be in `CamelCase`; unexported in `camelCase`.

## Testing Guidelines
- Use Go’s built-in `testing` package unless another framework is explicitly added.
- Name test files `*_test.go` and test functions `TestXxx`.
- Keep tests deterministic; fixed inputs are required for candidate ordering expectations.

## Commit & Pull Request Guidelines
- No commit message convention is established in this repository yet.
- Use short, imperative summaries (for example, “Add weighted node support”).
- PRs should include: a clear description of changes, rationale, and test evidence (command output or notes).

## Upstream Reference
- This repository references https://github.com/sile/rendezvous_hash. If you mirror or port functionality, note the upstream version or commit in the PR description.
