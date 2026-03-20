# Operator Runbook

## Audience

Use this runbook when building, smoke-testing, or troubleshooting an application that embeds `github.com/colt3k/mycli`.

## Build and Smoke Test

```bash
go build -mod=mod ./...
go test -mod=mod ./...
go run -mod=mod ./example -h
go run -mod=mod ./example -c example/config.toml server
```

If you changed dependencies, run `go mod vendor` before shipping so the committed `vendor/` tree matches `go.mod`.

## Routine Operations

### Show Help and Version

```bash
go run -mod=mod ./example -h
go run -mod=mod ./example -version
```

### Enable Environment Variables

The library ignores env vars unless the application sets `DisableEnvVars = false`. Remember that the default prefix is `T_`.

### Install Bash Completion

```bash
go install .
cp ./bash_autocomplete /usr/local/etc/bash_completion.d/mycli
```

Rename the completion file if the final executable name is not `mycli`.

## Troubleshooting

| Symptom | Likely Cause | Action |
| --- | --- | --- |
| `go build ./...` fails with inconsistent vendoring | `vendor/modules.txt` is stale | Run `go mod vendor`, or use `-mod=mod` while developing |
| `!!! no command set to run` | No command matched and `MainAction` is nil | Pass a valid command or configure `MainAction` |
| `required flag '-x' not set` | Final value still equals the default | Provide the flag on the command line, via env, or in config |
| Config value is ignored | Wrong TOML path or a command-line/env value already won | Check precedence and table names such as `[server]` or `[weserve.config]` |
| Env value is ignored | Env lookup disabled or wrong prefix | Set `DisableEnvVars = false` and verify `EnvPrefix` |
| Duplicate variable warning appears | Two flags share the same pointer with different defaults | Split the backing variables, or intentionally set `DisableFlagValidation = true` |

## Diagnostic Flags

- `-debug` enables debug logging.
- `-debugLevel` sets a more specific debug level for applications that honor it.
- `-generate-bash-completion` prints available completions instead of running the normal action.

## Recovery Steps

1. Re-run the command with `-h` to confirm the expected command and flag names.
2. Re-run with explicit command-line values to bypass env/config ambiguity.
3. Validate the config path and inspect the exact TOML table names.
4. Rebuild with `-mod=mod` if vendoring noise is blocking triage.
5. Sync `vendor/` with `go mod vendor` once the dependency set is correct.
