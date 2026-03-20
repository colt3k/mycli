# Dataflow

## Inputs

`mycli` merges four inputs into a single execution flow:

1. default values declared on flag structs
2. command-line arguments in `os.Args`
3. environment variables, when enabled
4. a TOML config file, when `-config` is supplied

## Parse Pipeline

```text
CLI definition
  -> add default flags
  -> optionally build environment variable names
  -> parse global flags once
  -> run PostGlblAction / VersionPrint
  -> reset the standard flag set
  -> rebuild global + command + subcommand flag sets
  -> overlay env values
  -> overlay config values
  -> validate required flags and Options
  -> handle help / version / bash completion
  -> resolve active command path
  -> run PreAction -> Action -> PostAction
```

## Internal State

- `FlgValues`: captures the first bound value for each flag key.
- `varMap`: records variable pointer reuse to warn about conflicting defaults.
- `TomlWrapper.Map`: stores parsed TOML as a nested map tree.
- `c.cur`: tracks the currently active command for help rendering.

## Config Data Path

`parseConfigFile()` walks three scopes in order:

1. global flags by key name
2. command flags by `command.flag`
3. subcommand flags by `command.subcommand.flag`

If a hidden command has a `Variable`, the entire subtree with the command name is unmarshaled into that value.

## Command Resolution

Command dispatch is positional. After global parsing, `Parse()` scans `os.Args` for a matching command name or short name. If a subcommand is found later in the argument list, that subcommand becomes the active command and its `FlagSet` parses the remaining arguments.

## Output Paths

- help text goes to stdout via `printUsage()` or command `FlagSet.Usage()`
- bash completion writes to `CLI.Writer`
- debug output is printed through `nglog`
- normal actions are provided entirely by the embedding application

## Special Modes

- `TestMode`: replaces exit-style flows with returns so tests can keep running.
- `ShowDuration`: prints nanosecond timings for major parse phases.
- `DisableEnvVars`: skips env naming and env overlay entirely.
- `DisableFlagValidation`: suppresses duplicate-pointer warnings during setup.
