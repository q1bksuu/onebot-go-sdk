---
name: go-test-gen
description: Generate Go unit tests, benchmarks, mocks, and coverage-driven test updates for this repository. Use when asked to create or update Go tests, improve coverage, generate benchmarks, or produce testify-based mocks following project conventions.
---

# Go Test Generator

## Understand the request

- Identify the target scope (file, package, or module) and mode (unit, coverage, benchmark, mock, update).
- Ask for missing inputs (file path, interface name, target coverage) before generating.

## Read project rules

- Read `go.mod` for the Go version and module path.
- Read `.golangci.yml` for lint-driven test conventions (paralleltest, thelper, testifylint).
- Read `codecov.yml` for coverage ignore paths and targets before coverage work.
- Read `CLAUDE.md` for project-specific testing rules.
- Read module docs (for example, `v11/entity/CLAUDE.md`) when behavior is unclear.

## Generate tests

- Prefer table-driven tests when 3+ scenarios exist.
- Use AAA (Arrange-Act-Assert) structure and clear sectioning.
- Use `t.Parallel()` where safe and required by lints.
- Use `t.Helper()` in helpers and `t.Cleanup()` for resource cleanup.
- Use testify `assert` / `require` for assertions.
- Name tests and subtests with descriptive snake_case.
- Skip unexported functions unless explicitly requested.

## Update existing tests

- Preserve existing structure, naming, and helper patterns.
- Add only the missing scenarios or branches.
- Avoid rewriting working tests unless necessary for correctness.

## Scenario heuristics

- Cover nil pointers, empty slices/maps, zero values, boundaries, and error branches.
- Exercise error paths (`if err != nil`) and typed error cases.
- Add concurrency tests for sync or channel usage with stable goroutine counts.
- Note generic functions that need manual type instantiation.

## Coverage-driven mode

- Run `go test -coverprofile` only when asked or when coverage is explicitly requested.
- Exclude files and directories ignored by `codecov.yml` when analyzing gaps.
- Focus on uncovered functions and branches; prioritize public APIs and edge cases.
- Report expected coverage deltas if you run coverage analysis.

## Benchmarks

- Use `b.ReportAllocs()` and `b.ResetTimer()`.
- Add `b.RunParallel()` when concurrency is relevant.

## Mocks

- Use `testify/mock` for interface mocks.
- Place mocks in a `mocks/` subfolder if one exists; otherwise keep in the same package with `_mock.go`.
- Provide helper constructors to set expectations and return values.

## Output

- Write tests in `*_test.go` next to the source file unless a separate package is requested.
- Run `gofmt` on generated files.
- Summarize generated tests and suggest next steps (run `go test`, update coverage).
