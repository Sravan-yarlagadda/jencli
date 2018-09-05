
[![Go Report Card](https://goreportcard.com/badge/github.com/Sravan-yarlagadda/jencli)](https://goreportcard.com/report/github.com/Sravan-yarlagadda/jencli)

# jencli - Jenkins Client 

### Jenkins client written in golang. This is capable of following tasks.

- start and monitor jobs


### Usage : 

```
jencli.exe --help
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

Use "jencli [command] --help" for more information about a command
```

#### Start command
```
λ jencli.exe start --help
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  jencli start [flags]

Flags:
  -h, --help           help for start
  -m, --monitor        Monitor job
  -t, --token string   token (default " ")
  -l, --url string     URL to start the job (default " ")
  -u, --user string    user (default " ")

Global Flags:
      --config string   config file (default is $HOME/.jencli.yaml)
```
      
### Sample Output
![alt text](https://github.com/Sravan-yarlagadda/jencli/blob/master/images/output.PNG)





