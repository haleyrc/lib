# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`github.com/haleyrc/lib` is a personal collection of purpose-built Go libraries. Each package is focused, minimal, and intended for export.

## Commands

- `make` — runs all checks (fmt, tidy, verify, vet, test)
- `make test` — `go test -v -count=1 -shuffle=on ./...`
- `make fmt` — `gofmt -s -w .`
- `make vet` — `go vet ./...`
- `make tidy` — `go mod tidy`
- `make verify` — `go mod verify`

Run a single test: `go test -v -count=1 -run TestName ./package/`

## Architecture

Nine independent packages, each single-responsibility:

- **assert** — test assertions returning a chainable `Result` type (`.Fatal()`, `.OK()`). Works with any `assert.T` implementation.
- **hash** — bcrypt password hashing. Implements `slog.LogValue()` to mask values in logs.
- **httputil** — request/response inspection (`RealIP`, `DumpRequest`, request ID propagation).
- **log** — structured JSON logging built on `slog`. Includes HTTP middleware (`ParseRequestInfo`) that decorates logs with request context.
- **markdown** — markdown with YAML frontmatter parsing/encoding via `Decoder`/`Encoder`.
- **router** — thin wrapper around `http.ServeMux` with method-specific routing.
- **server** — `http.Server` wrapper with sensible defaults (timeouts, graceful shutdown).
- **sqlite** — thin `sqlx.DB` wrapper for SQLite.
- **web** — HTTP response construction helpers (`JSON`, `StatusCode`, `Header`, etc.).

## Conventions

- **Functional options** for configuration (e.g., `log.Debug()`, `log.Output(w)`).
- **Error wrapping** with `fmt.Errorf("...: %w", err)`.
- **Example-based tests** (`ExampleXxx` functions) are the primary testing pattern. They serve as both tests and documentation.
- **Semantic commits**: `feat(package):`, `chore:`, `docs:`, `test:`.
- **Panic on unrecoverable init errors** (e.g., `hash.New`, `sqlite.Open`) — this is intentional.
- Minimal external dependencies.
