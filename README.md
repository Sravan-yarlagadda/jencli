# jencli - Jenkins Client 

### Jenkins client written in golang. This is capable of following tasks.

- start and monitor jobs


### Usage : 

``` jencli.exe --help
Jenkins commmand line interface. this can be used to perform following actions
        - Start Jenkins job and display the log.

Usage:
  jencli [command]

Available Commands:
  config      Configure username and token
  help        Help about any command
  start       Start a jenkins job

Flags:
      --config string   config file (default is $HOME/.jencli.yaml)
  -h, --help            help for jencli
  -t, --toggle          Help message for toggle

Use "jencli [command] --help" for more information about a command ```


