# API Reference

## Package `mycli`

### Construction

`NewCli(f FatalAdapter, u UsageAdapter) *CLI` creates the root CLI object. Pass `nil` for the default fatal/help adapters.

### Core Types

#### `CLI`

`CLI` is the application root. Important fields:

- `AppInfo`: title, version, description, author, and build metadata.
- `Flgs`: global flags.
- `Cmds`: top-level commands.
- `PostGlblAction`: hook that runs after global flag parsing.
- `MainAction`: fallback action when no command is matched.
- `DisableEnvVars`: disables env lookup when `true` (default).
- `EnvPrefix`: environment-variable prefix, default `"T"`.
- `DisableFlagValidation`: suppresses duplicate-pointer warnings.
- `ShowDuration`: prints timing for parse stages.
- `Writer`: destination for bash-completion output.
- `TestMode`: prevents exit-style flows during tests.

Common methods:

- `Parse() error`: parse input, apply overlays, and dispatch actions.
- `Help() bool`: reports whether top-level help was requested.
- `Command(name string) *CLICommand`: retrieves a top-level command.
- `Flag(name string, flgs []CLIFlag) CLIFlag`: finds a flag by name.
- `IsDebug() bool`, `DebugLevel() int64`: expose global debug state.
- `IsProxySet() bool`, `GetHttpProxy()`, `GetHttpsProxy()`, `GetNoProxy()`: expose proxy values.

#### `CLICommand`

`CLICommand` defines a command or subcommand:

- `Name`, `ShortName`, `Usage`
- `Flags`
- `SubCommands`
- `PreAction`, `Action`, `PostAction`
- `BashCompletion`
- `Hidden`
- `Variable`: used for hidden structured config payloads

`Action`, `PreAction`, and `PostAction` must be `func()` or `func() error`.

#### `CLIFlag`

`CLIFlag` is the interface implemented by every flag type. Implementations must support:

- flag binding to a `flag.FlagSet`
- environment lookup
- config lookup
- required-value checks
- option validation
- help rendering metadata

## Built-in Flag Types

- `BoolFlg`: boolean flags
- `Float64Flg`: `float64` flags
- `Int64Flg`: `int64` flags
- `StringFlg`: string flags, including comma-separated option validation
- `Uint64Flg`: `uint64` flags
- `VarFlg`: custom `flag.Value` wrapper using `StringList`

Each flag type accepts the same core fields: `Variable`, `Name`, `ShortName`, `Usage`, `Value`, `Required`, `Options`, `Hidden`, `EnvVar`, and `EnvVarExclude`.

## Config Types

- `Toml() *TomlWrapper`: returns the singleton TOML wrapper.
- `TomlWrapper`: loads a TOML file into a map and resolves dotted paths.
- `FixPath(path string) string`: converts relative paths to absolute paths before config loading.

## Package `custom`

`custom.TomlFlg` shows how to implement a structured flag that stores a `Clients` value. Supporting types:

- `Clients`
- `Host`
- `Connection`
- `Cert`

Use this package as the template when you need non-scalar config-backed flags.

## Errors

- `InvalidObjectError`: returned when a flag definition is not a pointer or is nil.
- `InvalidValueError`: returned when a flag value is outside the allowed `Options`.

## Minimal Example

```go
cli := mycli.NewCli(nil, nil)
cli.Title = "demo"
cli.DisableEnvVars = false

var port int64
cli.Cmds = []*mycli.CLICommand{
	{
		Name:   "serve",
		Usage:  "start server",
		Action: func() { fmt.Println("serving") },
		Flags: []mycli.CLIFlag{
			&mycli.Int64Flg{Variable: &port, Name: "port", ShortName: "p", Value: 8080},
		},
	},
}

if err := cli.Parse(); err != nil {
	log.Fatal(err)
}
```
