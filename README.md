# ecnu-cli

**This repo is forked from https://github.com/hduhelp/hdu-cli**

## Installation Or Upgrade

```shell
go install github.com/DarkHighness/ecnu-cli@latest
```

or direct download the release file which suffix match your platform.

## Startup

use command like

```
ecnu-cli net login --username {Your student number} --password {Your HDU Cas Password} --save
```

or manually use the .ecnu-cli.yaml and fill according the comments

<details>
<summary>Trouble shoot</summary>

> The Command may need root privilege
>
> and sometimes go env is not install completely on your root account (sudo mode)
>
> so try like `sudo $GOROOT/bin/go install github.com/DarkHighness/ecnu-cli@latest`
> 
> By the way, if you follow the offical installation guide of GO, The goroot will be like /usr/local/go/
</details>

## Usage

### ecnu-cli [sub command]

### Available Sub Commands:

- completion  
  - generate the autocompletion script for the specified shell
- help        
  - Help about any command
- net         
  - i-hdu network auth cli

### Flags:

- --config string   
  - config file (default is $HOME/.ecnu-cli.yaml)
  - more detail see comments at [ecnu-cli.yaml example](./.ecnu-cli.yaml)
- -h, --help            
  - help for ecnu-cli
- -s, --save            
  - save config
- -V, --verbose         
  - show more info
- -v, --version         
  - version for ecnu-cli


Use `ecnu-cli [sub command] --help` for more information about a command.


