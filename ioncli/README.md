# ioncli - CLI for Ionburst

## Installation

### Using Go

```sh
$ go install gitlab.com/ionburst/ionburst-sdk-go/ioncli
```

### Binaries

Please review the tagged releases page [HERE](https://gitlab.com/ionburst/ionburst-sdk-go/-/releases)

## Usage

```sh
$ ioncli -h
NAME:
   ioncli - Command Line Utility for Ionburst
USAGE:
   ioncli [global options] command [command options] [arguments...]
   
COMMANDS:
   classifications, class  Manage Ionburst Classifications
   get                     Download an object from Ionburst
   put                     Upload an object to Ionburst
   delete                  Delete an object from Ionburst
   help, h                 Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -v  Show debug/verbose output (default: false)
   --help, -h   show help (default: false)
   --version    Show Version Information (default: false)
   
COPYRIGHT:
   (C) Ionburst Limited UK
   
VERSION:
   V0.0.1 [ Build: a2e2f0b | API: 1.0.0 ]
```

### Get Classifications

```sh
$ ioncli class list
Classifications: 1

Classification
-----------------
Restricted

```

### Upload an object

```sh
$ ioncli put <OBJECT_ID> <FILENAME>

$ ioncli put testing_file testingfile

```

### Download an object

```sh
$ ioncli get <OBJECT_ID> <OUTPUT_FILENAME>

$ ioncli get testing_file testingfile

```

### Delete an object

```sh
$ ioncli get <OBJECT_ID> <OUTPUT_FILENAME>

$ ioncli get testing_file testingfile

```

## Building

After downloading the SDK source to your GOPATH

```sh
$ cd $GOPATH/src/gitlab.com/ionburst/ionburst-sdk-go 
## Or where it was downloaded

$ go mod tidy

## Compile bin/ioncli for your platform

$ make compile

## Cross compile for all platforms

$ make cross-compile

## Look in the bin folder for built binaries

```