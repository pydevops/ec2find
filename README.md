
The tool uses aws-sdk-go and is inspired by my coworker Erik Maciejewsk's ec2find, which is  written in Python and boto library.  

As a cloud engineer working with AWS, finding an instance IP private address and ssh into it could be a daily task. A tool like this would improve efficiency. My goal is to write something useful in `golang` as it can be compiled into a single binary, ideal for devops toolchain.

### build
`go install` and https://github.com/Masterminds/glide is used for managing package dependencies

### Usage

```
NAME:
   EC2 Instance Finder - Use this app to find the IP addresses of the ec2 instance by searching a given tag name

USAGE:
   ec2find [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --profile profile, -p profile                         aws profile (default: "default")
   --status instance-state-name, -s instance-state-name  instance-state-name (default: "running")
   --deploy_group, -d                                    search by deploy_group tag instead
   --login, -l                                           Prompt to log-in (ssh) after list is returned
   --help, -h                                            show help
   --version, -v                                         print the version
```

### Examples

```ec2find jenkins```
list instances uses AWS profile `default`.

```ec2find -p prod jenkins```
list instances uses AWS profile `prod`.

```ec2find -p prod -s stopped jenkins```
list stopped instances uses AWS profile `prod`.


```ec2find -l jenkins```
list instances uses AWS profile `default`, prompts ssh to the selected ec2 instance,  ec2find shells out the native ssh with the IP address found.
