# Repository Guidelines

## Project Structure & Module Organization

This repository is a small Go module centered on the root `mycli` package. Core CLI parsing, config loading, error handling, bash completion, and built-in flag types live at the repository root in files such as `cli.go`, `config.go`, `bashcompletion.go`, and `flg*.go`. Use `custom/` for extension flag types like [`custom/flgtoml.go`](custom/flgtoml.go). The runnable demo app is in `example/`, with sample input in `example/config.toml`. A `vendor/` tree is committed and should stay aligned with `go.mod`.

## Build, Test, and Development Commands

- `go build -mod=mod ./...` builds all packages while ignoring stale vendoring metadata.
- `go test -mod=mod ./...` runs the unit test suite.
- `go test -mod=mod -cover ./...` reports coverage for touched packages.
- `go mod vendor` refreshes `vendor/modules.txt` after dependency changes so plain `go build ./...` works again.
- `./buildAuto.sh` installs the package and copies `bash_autocomplete` into the local bash-completion directory.

## Coding Style & Naming Conventions

Use standard Go formatting with `gofmt`; tabs are the expected indentation. Keep exported identifiers in PascalCase and private helpers in lowerCamelCase. Follow existing API naming instead of “correcting” it: new flag implementations should match the current `*Flg` suffix pattern (`StringFlg`, `Int64Flg`, `VarFlg`). Keep packages focused: root for shared CLI behavior, `custom/` for optional extensions, and `example/` for runnable documentation.

## Testing Guidelines

Tests live beside source files as `*_test.go` and use the standard `testing` package with `github.com/stretchr/testify/assert`. Prefer table-driven tests like `TestFlgTypes` and `TestCmdHelp`. When changing flag parsing, environment-variable handling, or config precedence, add coverage for command-line, env, config, and default-value paths.

## Commit & Pull Request Guidelines

Recent history uses short, imperative, lowercase subjects such as `update deps` and `remove os.Exit calls`. Keep commits narrow and behavior-focused. Pull requests should summarize user-visible CLI changes, list the verification command(s) you ran, and call out updates to `README.md`, `example/`, `bash_autocomplete`, or `vendor/` when applicable.

## Configuration Notes

Flag resolution order is command line, then environment, then config file, then defaults. Preserve that precedence, and document any change that affects CLI compatibility.
