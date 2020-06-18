# MyCLI

  Abilities
  - toml configuration file
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
