# MyCLI

`mycli` is a Go library for building command-line applications with global flags, commands, subcommands, TOML-backed configuration, optional environment-variable loading, and bash completion.

## Documentation

- [Developer guide](docs/developer-guide.md)
- [API reference](docs/api-reference.md)
- [Config schema](docs/config-schema.md)
- [Dataflow](docs/dataflow.md)
- [Operator runbook](docs/operator-runbook.md)
- [Contributor guide](AGENTS.md)

## Abilities

### TOML configuration file

Pass `-config` to load values from TOML. Global flags live at the root, command flags live under `[command]`, and subcommand flags live under `[command.subcommand]`.

```toml
capture = "hello"

[server]
protocol = "https"
port = 9090

[weserve.config]
application = "gc"
```

Structured payloads also work. The sample in [`example/config.toml`](example/config.toml) uses `[[clients]]` to populate `custom.Clients`.

### Prefix to environment values

Environment lookup is disabled by default. Enable it with `cli.DisableEnvVars = false`. When enabled, `EnvPrefix` defaults to `"T"`, so `capture` maps to `T_CAPTURE`. Explicit `EnvVar` overrides are still prefixed unless you set `cli.EnvPrefix = ""`.

```go
cli := mycli.NewCli(nil, nil)
cli.DisableEnvVars = false
cli.EnvPrefix = "T"
```

### Commands and subcommands

```go
cli := mycli.NewCli(nil, nil)
cli.Title = "Example Title"
cli.Description = "Example Description"
cli.Version = version.VERSION

cli.Flgs = []mycli.CLIFlag{
	&mycli.StringFlg{Variable: &logDir, Name: "log_dir", ShortName: "ld", Usage: "override logging directory"},
}

cli.Cmds = []*mycli.CLICommand{
	{
		Name:   "update",
		Usage:  "check for updates",
		Action: func() error { return update() },
	},
	{
		Name:      "server",
		ShortName: "s",
		Usage:     "start server",
		Action:    func() error { return startServer(ctx, port) },
		Flags: []mycli.CLIFlag{
			&mycli.Int64Flg{Variable: &port, Name: "port", ShortName: "p", Value: 8081},
		},
	},
}

err := cli.Parse()
```

Subcommands are nested on `CLICommand.SubCommands`, as shown in [`example/main.go`](example/main.go) for `weserve config` and `weserve cmdln`.

### Global and command flags

Global flags belong in `cli.Flgs`. Command-local flags belong in `CLICommand.Flags`. `Parse()` also injects built-in flags for help, debug, debug level, version, config, proxy values, and bash completion when applicable.

### Custom and default flag types

Built-in flag types:

- `BoolFlg`
- `Float64Flg`
- `Int64Flg`
- `StringFlg`
- `Uint64Flg`
- `VarFlg` (`StringList`)

Custom flag type included in this repo:

- `custom.TomlFlg`

### Flag attributes

Required flags:

```go
&mycli.Int64Flg{Variable: &port, Name: "port", ShortName: "p", Usage: "Set Port", Value: 9111, Required: true}
```

Limited option sets:

```go
&mycli.StringFlg{
	Variable: &capture,
	Name:     "capture",
	ShortName:"cap",
	Usage:    "Used to test string",
	Options:  []string{"hello", "bye"},
}
```

### Help

`-h` prints global usage, commands, subcommands, defaults, and option metadata. Command help is also available on individual commands, for example `server -h`.

### Bash autocompletion

Use the included `bash_autocomplete` script with bash-completion v2+.

```bash
go install .
cp ./bash_autocomplete /usr/local/etc/bash_completion.d/mycli
```

Rename the installed completion file to match your executable if you embed this library in another application.

## Order of precedence on flag values

1. Command line
2. Environment variables
3. Config file
4. Defaults

Environment values only participate when `DisableEnvVars` is set to `false`.

## Testing using the example

Run the demo application from the repository root:

```bash
go run -mod=mod ./example -c example/config.toml server
go run -mod=mod ./example -c example/config.toml client
go run -mod=mod ./example -c example/config.toml weserve cmdln
go run -mod=mod ./example -c example/config.toml weserve config
```

## Warning on reuse of variables across commands

If the same variable pointer is reused across multiple flags with different defaults, the library prints a warning because later bindings can override earlier values unexpectedly. Disable this validation only when the overlap is intentional:

```go
cli.DisableFlagValidation = true
```
