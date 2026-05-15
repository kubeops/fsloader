# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `kubeops.dev/fsloader` — a small sidecar/CLI that **watches a file and runs a bash script** when the file changes. Used to decouple application config files from container images: mount a `ConfigMap`/`Secret` into a pod, point `fsloader` at the mounted file, and it'll re-run a reload script on every change (or write the new contents to a different destination and reload).

The produced binary is `fsloader`. Builds for **5 platforms** (linux amd64/arm/arm64 + windows/amd64 + darwin/amd64 + darwin/arm64) — sometimes used outside Kubernetes too.

## Architecture

- `main.go` — entry point at the module root; wires `cmds.NewRootCmd()` and runs.
- `cmds/`:
  - `root.go` — Cobra root command.
  - `run.go` — the long-running watcher. Uses `fsnotify` (vendored) to watch the source file; on change writes to the destination and invokes the configured shell command.
- `version.go` — version variables injected via `-ldflags`.
- `docs/reference/` — generated CLI reference.
- `hack/`, `Makefile` — AppsCode build harness.
- `vendor/` — checked-in deps.

This is intentionally a tiny binary; resist adding controller-runtime, codegen, or CRDs.

## Common commands

All Make targets run inside `ghcr.io/appscode/golang-dev` — Docker must be running.

- `make ci` — CI pipeline.
- `make build` — build for the host OS/ARCH into `bin/<os>_<arch>/fsloader`.
- `make all-build` — build for every `BIN_PLATFORMS` (linux amd64/arm/arm64 + windows/amd64 + darwin/amd64 + darwin/arm64).
- `make fmt`, `make lint`, `make unit-tests` / `make test` — standard.
- `make verify` — `verify-gen verify-modules`; `go mod tidy && go mod vendor` must leave the tree clean.
- `make add-license` / `make check-license` — manage license headers.

No container target — this binary is typically embedded in an upstream sidecar image rather than published standalone.

Run a single Go test (requires a local Go toolchain):

```
go test ./... -run TestName -v
```

## Conventions

- Module path is `kubeops.dev/fsloader` (vanity URL). Imports must use that.
- License: Apache-2.0 (`LICENSE`); new files need the standard "Copyright AppsCode Inc. and Contributors" header (`make add-license`).
- Sign off commits (`git commit -s`); contributions follow the DCO (`DCO`, `CONTRIBUTING.md`).
- Vendor directory is checked in — `go mod tidy && go mod vendor` must leave the tree clean (enforced by `verify-modules`).
- Keep this binary tiny — single-purpose "watch file, run script". If a use case needs templating or multi-file orchestration, build a separate binary; don't graft features here.
- Builds linux/windows/darwin host binaries; do not pull in linux-only or cgo deps.
- `main.go` lives at the module root (not under `cmd/`); keep that convention.
