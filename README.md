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

            Example
            c := mycli.NewCli(nil, nil)
            c.EnvPrefix = "T"
  - Commands

            Example
            const (
              title = "Example Title"
              description = "Example Description"
              author = "John Doe (j.doe@example.com)"
              copyRight = "(c) YYYY Example Inc.,"
            )

            func setupFlags(ctx context.Context) {
              var overrideLogDir string
              var port int64 = 8081
              cli := mycli.NewCli(nil, nil)

              cli.Title = title
              cli.Description = description
              cli.Version = version.VERSION
              cli.BuildDate = version.BUILDDATE
              cli.GitCommit = version.GITCOMMIT
              cli.GoVersion = version.GOVERSION
              cli.Author = author
              cli.Copyright = copyRight
              cli.PostGlblAction = func() error { return setLogger(overrideLogDir) }
              cli.Flgs = []mycli.CLIFlag{
		         &mycli.StringFlg{Variable: &overrideLogDir, Name: "log_dir", ShortName: "ld", Usage: "override logging directory"},
              }

              cli.Cmds = []*mycli.CLICommand{
                  {
                      Name:   "update",
                      Usage:  "check for updates",
                      Action: func() error { return update() },
                      Flags: []mycli.CLIFlag{},
                  },
                  {
			       Name:   "serve",
			       Usage:  "start server",
                       Action: func() error { return startServer(ctx, port) },
			       Flags: []mycli.CLIFlag{
				   &mycli.Int64Flg{Variable: &port, Name: "port", ShortName: "p", Value: 8081},
                       },
                  },
              }
	          err := cli.Parse()
	          if err != nil {
		          log.Logf(log.FATAL, "%v", err)
	          }
           }
  - Sub Commands

            EXAMPLE
            cli := mycli.NewCli(nil, nil)
            ...
            cli.Cmds = []*mycli.CLICommand{
    		Name:  "client",
    		Usage: "use as a client",
    		SubCommands: []*mycli.CLICommand{
    			{
    				Name:      "config",
    				ShortName: "c",
    				Usage:     "use config file",
    				Action: func() { log.Println("ran clients config") },
    				Flags: []mycli.CLIFlag{
    					&mycli.Int64Flg{Variable: &t4, Name: "port", ShortName: "p", Usage: "Set Port", Value: 9111, Required: false},
    				},
    			},
    			{
    				Name:      "cmdln",
    				ShortName: "cl",
    				Usage:     "use command line",
    				Action: func() { log.Println("ran clients cmdline") },
    				Flags: []mycli.CLIFlag{
    					&mycli.Int64Flg{Variable: &t5, Name: "port", ShortName: "p", Usage: "Set Port", Value: 9111, Required: false},
    				},
    			},
    		},
     	    }
  - Global and Command Flags
  
           EXAMPLE
           cli := mycli.NewCli(nil, nil)
           ...
           cli.Flgs = []mycli.CLIFlag{
    	      &mycli.StringFlg{Variable: &overrideLogDir, Name: "log_dir", ShortName: "ld", Usage: "override logging directory"},
           }
  - Custom Flag types 
    - toml
  - Default Flag types
    - bool
    - float64
    - int64
    - string
    - uint64
    - toml
  - Help
          
         Example
         $ main -h
         NAME:
           main

         USAGE:
           main [global options] command [command options] [arguments...]

         GLOBAL OPTIONS:
           -debug, -d
                 flag set to debug (default false)

           -debugLevel, -dbglvl  int
                 set debug level (default 0)

           -config, -c  string
                 config file path

           -proxyhttp  string
               HTTP_PROXY  (as environment var)
                 Sets http_proxy for network connections

           -proxyhttps  string
               HTTPS_PROXY (as environment var)
                 Sets https_proxy for network connections

           -noproxy  string
               NO_PROXY    (as environment var)
                 Sets no_proxy for network connections

           -log_dir, -ld  string
                 override logging directory

         COMMANDS:
           update:    (check for updates)

           serve:     (start server)
               -port, -p  int
                   Set port (default 8081)

           client:    (use as a client)

             Sub Commands:
               config :  use config file
                 -port, -p  int
                    Set Port (default 9111)

               cmdln :   use command line
                -port, -p  int
                    Set Port (default 9111)


  - Flag Attributes
    - required flag

          EXAMPLE
          &mycli.Int64Flg{Variable: &t4, Name: "port", ShortName: "p", Usage: "Set Port", Value: 9111, Required: true},
    - option values, limit to set of valid options
  
          EXAMPLE
          &mycli.StringFlg{Variable: &capture, Name: "capture", ShortName: "cap", Usage: "Used to test string", Options: []string{"hello", "bye"}},
  - Bash Autocompletion
    #### use the included bash_autocomplete script along with bash-completion v2+
        Typically you will perform the following to set this up
        - Name of executable i.e. mytest
        - Copy bash_autocomplete to /etc/bash_completion.d/ using the name of the executable 
           i.e. /etc/bash_completion.d/mytest
        - chmod 777 for all to use i.e. chmod 777 /etc/bash_completion.d/mytest
    
Order of precedence on FLAG values
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