// This program runs a program at the time specified.
//
// It is used like this:
//    runat 12:31:22 my prog -a -b
//
// I wrote it to allow me to analyze process concurrency issues.
//
// Copyright (c) 2016 by Joe Linoff
// Licence: MIT Open Source
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var version = "v0.2" // 2016-07-22

func main() {
	ts, cmd, v := getOptions()
	if v > 0 {
		Info("timestamp: %v", ts)
		Info("command  : %v", getCmdString(cmd))
	}

	// Get the start time.
	start := getTimeToStart(ts, v)
	if v > 0 {
		Info("start    : %v", start)
	}

	// Wait until it is time to start and
	// then launch the process.
	ltime := wait(start, v)
	if v > 0 {
		Info("launchat : %v", ltime)
		Info("launch   : (%v) %v", len(cmd), getCmdString(cmd))
	}

	// Launch the process.
	launch(cmd, v)
}

// Launch the command.
func launch(cmd []string, v int) {
	env := os.Environ()
	binary, err := exec.LookPath(cmd[0])
	if err != nil {
		Err("%v: '%v'", err, cmd[0])
	}
	binary, err = filepath.Abs(binary)
	if err != nil {
		Err("%v: '%v'", err, binary)
	}
	if v > 0 {
		Info("binary   : %v", binary)
	}
	if _, err = os.Stat(binary); os.IsNotExist(err) {
		Err("program does not exist: '%v'", cmd[0])
	}
	err = syscall.Exec(binary, cmd, env)
	if err != nil {
		Err("%v: '%v'", err, binary)
	}
}

// Get the command string.
// It will quote arguments with spaces or quotations.
// Here are some examples:
//    pwd                 --> pwd
//    echo foo            --> echo foo
//    echo "foo bar"      --> echo "arg1a arg1b"
//    echo "it's great!"  --> echo "it's great!"
//    echo '"quote"'      --> echo '"quote"'
func getCmdString(cmd []string) string {
	cs := ""
	re1 := regexp.MustCompile(`[\"\' \t]`)
	for i, arg := range cmd {
		if i > 0 {
			cs += " "
		}
		if re1.MatchString(arg) {
			// Contains whitespace or quotes.
			if strings.Contains(arg, `"`) == false {
				// Easy!
				cs += `"`
				cs += arg
				cs += `"`
			} else if strings.Contains(arg, `'`) == false {
				// Easy!
				cs += `'`
				cs += arg
				cs += `'`
			} else {
				// Contains both. Assume a single quote but that may not always be
				// correct. For this application it really doesn't matter.
				p := ""
				cs += `'`
				for i := 0; i < len(arg); i++ {
					c := string(arg[i])
					if c == `'` && p != "\\" {
						cs += "\\"
					}
					cs += c
					p = c
				}
				cs += `'`
			}
		} else {
			cs += arg
		}
	}
	return cs
}

// Wait to start by polling the current time.
func wait(start time.Time, v int) time.Time {
	t := time.Now()
	d := start.Sub(t)
	if d.Seconds() > 1.1 {
		// Sleep for a bit before doing detailed polling
		// checks. Make sure to truncate so that we have
		// sufficient margin for the precision.
		s := time.Duration(time.Second * time.Duration(int(d.Seconds())))
		if v > 1 {
			Info("sleep    : %v", s)
		}
		time.Sleep(s)
	}
	if v > 1 {
		Info("polling")
	}
	for {
		if t.Equal(start) || t.After(start) {
			break
		}
		t = time.Now()
	}
	return t
}

// Get the time to start.
func getTimeToStart(ts string, verbose int) time.Time {
	// First verify that the time specification has the proper format.
	hr := -1
	min := -1
	sec := -1
	re1 := regexp.MustCompile(`^(\d+):(\d+):(\d+)$`)
	re2 := regexp.MustCompile(`^(\d+)$`)
	if re1.MatchString(ts) {
		group := re1.FindAllStringSubmatch(ts, -1)
		hr, _ = strconv.Atoi(group[0][1])
		min, _ = strconv.Atoi(group[0][2])
		sec, _ = strconv.Atoi(group[0][3])
		if hr < 0 || hr > 23 {
			Err("time specification '%v' has invalid hour: %v, must be in the range [0...23]", ts, hr)
		}
		if min < 0 || min > 69 {
			Err("time specification '%v' has invalid minute: %v, must be in the range [0...60]", ts, min)
		}
		if sec < 0 || sec > 69 {
			Err("time specification '%v' has invalid second: %v, must be in the range [0...60]", ts, sec)
		}
	} else if re2.MatchString(ts) {
		group := re2.FindAllStringSubmatch(ts, -1)
		sec, _ = strconv.Atoi(group[0][1])
		if sec > 59 {
			Err("time specification %v seconds out of range [0..59]", sec)
		}
	} else {
		Err("unrecognized time specification: '%v', see help (-h) for more information", ts)
	}

	// Figure out the start time.
	then := time.Now()
	if verbose > 1 {
		Info("check    : %v", then)
	}

	if sec >= 0 {
		if sec < then.Second() {
			// Advance to the next minute if the current time is
			// later than the specified time.
			// Example:
			//   current time: 45
			//   specified time: 20
			//   advance 15 (60-45) seconds.
			d := time.Duration(time.Second * time.Duration(60-then.Second()))
			then = then.Add(d)
			if verbose > 1 {
				Info("added    : %v (%v)", d, sec)
				Info("check    : %v", then)
			}
		}

		// Now set the correct second offset.
		if sec > then.Second() {
			// The specified time is later than the current time.
			// This is easy, set the explicit time.
			d := time.Duration(time.Second * time.Duration(sec-then.Second()))
			then = then.Add(d)
			if verbose > 1 {
				Info("added    : %v (%v)", d, sec)
				Info("check    : %v", then)
			}
		}
	}

	// Set the minute explicitly.
	if min >= 0 {
		if min < then.Minute() {
			// Advance to the next hour.
			d := time.Duration(time.Minute * time.Duration(60-then.Minute()))
			then = then.Add(d)
			if verbose > 1 {
				Info("added    : %v (%v)", d, sec)
				Info("check    : %v", then)
			}
		}

		if min > then.Minute() {
			d := time.Duration(time.Minute * time.Duration(min-then.Minute()))
			then = then.Add(d)
			if verbose > 1 {
				Info("added    : %v (%v)", d, sec)
				Info("check    : %v", then)
			}
		}
	}

	// Set the hour explicitly.
	if hr >= 0 {
		if hr < then.Hour() {
			// Advance to the next day (24 hour clock).
			d := time.Duration(time.Hour * time.Duration(24-then.Hour()))
			then = then.Add(d)
			if verbose > 1 {
				Info("added    : %v (%v)", d, sec)
				Info("check    : %v", then)
			}
		}

		if hr > then.Hour() {
			d := time.Duration(time.Hour * time.Duration(hr-then.Hour()))
			then = then.Add(d)
			if verbose > 1 {
				Info("added    : %v (%v)", d, sec)
				Info("check    : %v", then)
			}
		}
	}

	// Get rid of the nanoseconds.
	tf := time.Date(then.Year(),
		then.Month(),
		then.Day(),
		then.Hour(),
		then.Minute(),
		then.Second(),
		0,
		then.Location())
	return tf
}

func getOptions() (string, []string, int) {
	// Options are allowed before the time specification.
	verbose := 0
	ts := ""
	var cmd []string
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-h", "--help":
			help()
			os.Exit(0)
		case "-v", "--verbose":
			verbose++
		case "-vv":
			verbose += 2
		case "-V", "--version":
			base := filepath.Base(os.Args[0])
			fmt.Printf("%v %v\n", base, version)
			os.Exit(0)
		default:
			ts = os.Args[i]
			j := i + 1
			if j >= len(os.Args) {
				Err("missing command, see help (-h) for more information (%v >= %v)", j, len(os.Args))
			}
			cmd = os.Args[j:]
			i = len(os.Args) // exit the loop
		}
	}
	if ts == "" {
		Err("time specification missing, see help (-h) for more information")
	}
	if len(cmd) == 0 {
		Err("command missing, see help (-h) for more information")
	}

	return ts, cmd, verbose
}

func help() {
	base := filepath.Base(os.Args[0])
	msg := `
USAGE
    %[1]v [OPTIONS] <timespec> <command>

DESCRIPTION
    Run a command at the specific time with a resolution of 1 second.

    This tool can be used for analyzing process race conditions because it
    allows you to start multiple processes at almost the same time (usually
    to within microseconds of one another) in a convenient way (without having
    to specify a specific time).

    There are two input formats:

      HH:MM:SS  Run the command at this specific time.

      MARK      Run command at the specific second mark from the current
                minute. If that time has passed, it will rollover to the mark
                in the next minute

    The mark format is the most interesting because you don't have to think
    about the specific time. If you specify 30, it will run the command at the
    30 second mark, if you specify 15 it will run at the 15 second mark. If the
    mark has already passed it will start it in the next minute. Some examples
    will make that clear.

      Example 1, if the current time is 14:22:17 and you specify 30,
                 it will start at 14:22:30.

      Example 2, if the current time is 14:22:17 and you specify 10, it will
                 start at 14:23:30.

    The command is any command. Everything after the time specification is part
    of the command.

    Here is an example usage using three windows that will run some commands at
    the 30 second mark (about 28 seconds after they were started).

      $ # What is the current time seconds?
      $ date +'%%S'
      12

      $ # Window 1
      $ %[1]v 30 /bin/bash -c "echo win1 && date && pwd"

      $ # Window 2
      $ %[1]v 30 /bin/bash -c "echo win2 && date && pwd"

      $ # Window 3
      $ %[1]v 30 /bin/bash -c "echo win3 && date && pwd"

    Output at the 30 second mark.

      # Window 1:
      win1
      Thu Jul 21 18:49:30 PDT 2016
      /Users/jlinoff/work/runat.work

      # Window 2:
      win2
      Thu Jul 21 18:49:30 PDT 2016
      /Users/jlinoff/work/runat.work

      # Window 3:
      win3
      Thu Jul 21 18:49:30 PDT 2016
      /Users/jlinoff/work/runat.work

    As you can see, they all started at the same second.

OPTIONS
    All of the options must appear before the time specification.

    -h, --help         Print this help message and exits.
    -v, --verbose      Increase the level of verbosity.
    -V, --version      Print the program version and exits.

EXIT STATUS
    Returns the exit status of the command unless the command line is not
    syntactically correct in which case it reports an error message and exits
    with status 1.

EXAMPLES
    $ # Example 1: help
    $ %[1]v -h

    $ # Example 2: not help, -h occurs after the time specification
    $ %[1]v 30 sleep 10 -h

    $ # Example 3: run a command at the 30 second mark
    $ %[1]v 30 sleep 5

    $ # Example 4: run a command at a specific time (01:47:10 PM today)
    $ #            use "cron" or "at" if you want more control
    $ %[1]v 13:47:10 sleep 5

    $ # Example 5: run two commands at the 12 second mark
    $ %[1]v -v 12 /bin/bash -c "date && pwd"
`
	msg += "\n"
	fmt.Printf(msg, base)
}

// Info reports an informational message to stdout.
// Called just like fmt.Printf.
func Info(f string, a ...interface{}) {
	Base(os.Stdout, "INFO", fmt.Sprintf(f, a...), 2, false)
}

// Warn reports a warning message to stdout.
// Called just like fmt.Printf.
func Warn(f string, a ...interface{}) {
	Base(os.Stdout, "WARNING", fmt.Sprintf(f, a...), 2, false)
}

// Err reports an Error message to stderr and exits.
// Called just like fmt.Printf.
func Err(f string, a ...interface{}) {
	Base(os.Stderr, "ERROR", fmt.Sprintf(f, a...), 2, false)
	os.Exit(1)
}

// ErrWithLevel reports an Error message to stderr and exits.
// Called just like fmt.Printf.
// level should be 3 to get the caller of the caller.
func ErrWithLevel(level int, f string, a ...interface{}) {
	Base(os.Stderr, "ERROR", fmt.Sprintf(f, a...), level, false)
	os.Exit(1)
}

// Base is the basis for the messages.
func Base(fp *os.File, p string, s string, level int, sf bool) {
	pc, fname, lineno, _ := runtime.Caller(level)
	fname = fname[0 : len(fname)-3]
	//t := time.Now().UTC().Truncate(time.Millisecond).String()
	t := time.Now().UTC().Truncate(time.Millisecond).Format("2006-01-02 15:05:05.000 MST")
	if sf { // show function name
		fct := runtime.FuncForPC(pc).Name()
		fmt.Fprintf(fp, "%-28s %-7s %s %s %d - %s\n", t, p, path.Base(fname), fct, lineno, s)
	} else {
		fmt.Fprintf(fp, "%-28s %-7s %s %4d - %s\n", t, p, path.Base(fname), lineno, s)
	}
}
