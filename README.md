# MyCLI

  Abilities
  - toml configuration file

    - allows you to configure all flags for all commands but you still  need to pass the command/subcommand to run
    - see example/config.toml for example
      
            Example
            globalflagName="XX"
      
            [command]
                [command.subcommand]
                flagname="XX"
  - prefix to environment values
  - Commands
  - Sub Commands
  - Global and Command Flags
  - Custom Flag types 
    - toml
  - Default Flag types
    - bool
    - float64
    - int64
    - string
    - uint64
    - toml
  - Flag Attriburtes
    - required flag
    - option values, limit to set of valid options
  - Bash Autocompletion
    
Order of precedence on values
   - 1. commandline  (highest priority)
   - 2. environment
   - 3. config file
   - 4. defaults     (lowest priority)
