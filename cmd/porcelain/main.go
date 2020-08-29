package main

import (
	"flag"
	"fmt"
	"github.com/robertgzr/porcelain"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// TODO allow custom log location
const logLoc string = "/tmp/porcelain.log"

func init() {
	flag.BoolVar(&porcelain.DebugFlag, "debug", false, "write logs to file ("+logLoc+")")
	flag.BoolVar(&porcelain.FmtFlag, "fmt", true, "print formatted output (default)")
	flag.BoolVar(&porcelain.BashFmtFlag, "bash", false, "escape fmt output for bash")
	flag.BoolVar(&porcelain.NoColorFlag, "no-color", false, "print formatted output without color codes")
	flag.BoolVar(&porcelain.ZshFmtFlag, "zsh", false, "escape fmt output for zsh")
	flag.BoolVar(&porcelain.TmuxFmtFlag, "tmux", false, "escape fmt output for tmux")
	flag.StringVar(&porcelain.Cwd, "path", "", "show output for path instead of the working directory")
	flag.BoolVar(&porcelain.VersionFlag, "version", false, "print version and exit")

	logToStderr := flag.Bool("logToStderr", false, "write logs to stderr")
	flag.Parse()

	if porcelain.VersionFlag {
		fmt.Printf("porcelain version %s (%s)\nbuilt %s\n", porcelain.Version, porcelain.Commit, porcelain.Date)
		os.Exit(0)
	}

	if porcelain.DebugFlag {
		var (
			err   error
			logFd io.Writer
		)
		if *logToStderr {
			logFd = os.Stderr
		} else {
			logFd, err = os.OpenFile(logLoc, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
			if err != nil {
				os.Exit(1)
			}
		}
		log.SetOutput(logFd)
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	if porcelain.Cwd == "" {
		porcelain.Cwd, _ = os.Getwd()
	}
}

func main() {
	log.Println("running porcelain...")
	log.Println("in directory:", porcelain.Cwd)

	var out string
	switch {
	case porcelain.FmtFlag:
		out = porcelain.Run().Fmt()
	default:
		flag.Usage()
		fmt.Println("\nOutside of a repository there will be no output.")
		os.Exit(1)
	}

	_, _ = fmt.Fprint(os.Stdout, out)
}
