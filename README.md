# runat
Start multiple processes in parallel at the same time for testing race conditions

### Installation
Here are the steps to install it locally. Make sure that you have a recent version of go available.

```bash
$ $ go version
go version go1.6.2 darwin/amd64
$ git clone https://github.com/jlinoff/runat.git
$ cd runat
$ GOBIN=$(pwd)/bin go install runat.go
$ bin/runat -h
```
### Using it.
This tool can be used for analyzing process race conditions because it allows you to start multiple
processes at almost the same time (usually to within microseconds of one another) in a convenient
way (without having to specify a specific time).

It has two sets of arguments, a time specification and the command you want to run. Running it looks like this:

```bash
$ runat 25 echo "hello, world"
```

The time specification has two input formats.

| Format | Description |
| :---   | :---------- |
| HH:MM:SS | Start the command at this specific time. |
| _mark_ | Run command at the specific second mark from the current minute. If that time has passed, it will rollover to the mark in the next minute|

The _mark_ format is the most interesting because you don't have to think about the specific time.
If you specify 30, it will run the command at the 30 second mark, if you specify 15 it will run at the 15 second mark.
If _mark_ has already passed it will start it in the next minute. Some examples will make that clear.

Example 1, if the current time is 14:22:17 and you specify 30, it will start at 14:22:30.

Example 2, if the current time is 14:22:17 and you specify 10, it will start at 14:**23**:30.

The command is any command. Everything after the time specification is part of the command.

### An Example
This example shows how to use the tool to start processes in two different windows.

```bash
$ date +'%S'
12

$ # Window 1:
$ runat 30 /bin/bash -c "echo win1 && date && pwd"

$ # Window 2:
$ 30 /bin/bash -c "echo win2 && date && pwd"
```

Just after the 30 second mark this is what you will see.

```bash
$ date +'%S'
30

$ # Window 1:
win1
Thu Jul 21 18:49:30 PDT 2016
/Users/jlinoff/work/runat.work

$ # Window 2:
win2
Thu Jul 21 18:49:30 PDT 2016
/Users/jlinoff/work/runat.work
```

### Options
The command has 3 options. They must be specified before the time specification so that they will not be confused with command options.

| Short | Long      | Description |
| :---- | :-------- | :---------- |
| -h    | --help    | Print the help message and exit. |
| -v    | --verbose | Increase the level of verbosity. Can specify -vv as a short cut for -v -v. |
| -V    | --version | Print the version number and exit. |

Please send comments to improve this tool if you find it useful.
