# Config Schema

## Overview

Configuration is loaded only when the built-in `-config` flag is passed. The file is parsed as TOML, relative paths are normalized to absolute paths, and values are applied only when a flag still holds its default.

## Resolution Order

1. Command-line arguments
2. Environment variables
3. Config file values
4. Default values

Environment variables participate only when `DisableEnvVars` is `false`.

## TOML Layout

### Global Flags

Root-level keys map to global flags:

```toml
capture = "hello"
path = "/tmp/data"
url = "https://example.com"
```

### Command Flags

Command-local flags map to tables named after the command:

```toml
[server]
protocol = "https"
port = 9090
```

### Subcommand Flags

Subcommand flags use nested tables:

```toml
[weserve.config]
application = "gc"
port = 9111
```

Array-of-table input is also accepted because the TOML walker resolves the last array element while traversing nested paths. The sample config uses this form for `weserve` and `clients`.

### Hidden Command Payloads

Hidden commands can bind entire TOML sections into structured types. In the example app, the hidden `clients` command is populated from an array of tables:

```toml
[[clients]]
name = "host1"

[clients.connection]
protocol = "ssh"
host = "8.8.8.8"
port = 22

[clients.cert]
certpath = "/some/path/to/cert"
```

## Environment Variable Naming

When env lookup is enabled, names are derived like this:

- default prefix: `T`
- default mapping: `capture` -> `T_CAPTURE`
- explicit name with prefix: `EnvVar: "TESTTWO"` + `EnvPrefix: "TST"` -> `TST_TESTTWO`
- no prefix: set `EnvPrefix = ""`

Built-in config and proxy flags follow the same rule. To use raw names like `HTTP_PROXY`, set `EnvPrefix = ""`.

## Supported Value Shapes

- scalars: `bool`, `float64`, `int64`, `string`, `uint64`
- comma-separated string lists via `VarFlg`
- structured config via custom flag implementations such as `custom.TomlFlg`

## Validation Rules

- `Required: true` means the final resolved value must differ from the default.
- `Options` restrict the accepted final value after command-line, env, and config overlays are applied.
- Duplicate variable pointers across flags produce a warning unless `DisableFlagValidation` is `true`.
