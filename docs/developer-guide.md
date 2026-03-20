# Developer Guide

## Scope

This repository contains a reusable CLI library in the root `mycli` package, a `custom` package with an example structured flag type, and an `example` application that demonstrates real wiring.

## Repository Layout

- `cli.go`: core parse lifecycle, command dispatch, help rendering, and default flag injection.
- `config.go`: TOML singleton wrapper and key-path lookup.
- `flags.go`, `flg*.go`: `CLIFlag` contract plus built-in flag implementations.
- `bashcompletion.go`: main and subcommand completion emitters.
- `custom/flgtoml.go`: example of a custom structured flag backed by TOML/JSON data.
- `example/`: runnable demo app and sample config.
- `cli_test.go`: integration-style tests for initialization, help, command dispatch, and flag behavior.

## Local Development

Use module mode for day-to-day work because the committed `vendor/` tree may lag `go.mod`:

```bash
go build -mod=mod ./...
go test -mod=mod ./...
go test -mod=mod -cover ./...
go mod vendor
```

Run `go mod vendor` after dependency changes so plain `go build ./...` matches the checked-in vendor metadata again.

## Core Execution Model

`NewCli()` constructs a `CLI` with default settings: environment loading disabled, `EnvPrefix` set to `"T"`, and top-level bash completion enabled. Callers then fill `AppInfo`, `Flgs`, `Cmds`, optional hooks, and finally call `Parse()`.

`Parse()` does the following:

1. Injects default flags (`help`, `debug`, `debugLevel`, `version`, `config`, proxy flags, and bash completion).
2. Builds initial global flags so built-ins can be parsed early.
3. Runs global env lookup and `PostGlblAction`.
4. Rebuilds the flag sets for globals, commands, and subcommands.
5. Overlays environment and config values onto any flag still at its default.
6. Validates required flags and option lists.
7. Resolves the active command/subcommand and runs `PreAction`, `Action`, and `PostAction`.

Only `func()` and `func() error` actions are supported.

## Extending the Library

To add a new flag type, follow the existing `flg*.go` pattern:

1. Add a struct that implements `CLIFlag`.
2. Bind a pointer-backed variable in `BuildFlag`.
3. Implement env/config retrieval, required checks, option validation, and `ValueAsString`.
4. Add tests in `cli_test.go` or a new `*_test.go` file.
5. Update `README.md` and the docs in `docs/` if the public behavior changes.

`custom.TomlFlg` is the best reference for a non-scalar implementation.

## Testing Notes

The test suite uses stubs for fatal/help adapters and relies on `ResetForTesting(nil)` to rebuild the standard `flag.CommandLine` state between tests. When changing parse behavior, cover all four value sources: defaults, config, environment, and explicit command-line arguments.

If you add new help text, completion output, or config semantics, update the example app and docs in the same change.

## Documentation Maintenance

Keep `README.md` focused on capabilities and quick-start usage. Put deeper implementation details in the docs under `docs/`. Do not add `doc.go` files in this repository; use Markdown and source comments instead.
