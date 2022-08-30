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

Testing using example
    - 1. go run main.go -c config.toml server
    - 2. go run main.go -c config.toml client
    - 3. go run main.go -c config.toml weserve cmdln
    - 4. go run main.go -c config.toml weserve config

Warning on Reuse of Variables across Commands
- When a variable is reused on different commands a warning is displayed
  - this will help inform you that the value could be overridden unexpectedly
  DISABLE by setting in your declaration
  - cli.DisableFlagValidation = true