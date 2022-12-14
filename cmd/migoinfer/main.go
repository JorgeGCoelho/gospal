// Command migoinfer is the command line entry point to MiGo type inference.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/JorgeGCoelho/gospal/v2/migoinfer"
	"github.com/JorgeGCoelho/gospal/v2/ssa/build"
)

const (
	Usage = `migoinfer is a tool for infering MiGo types from Go source code.

Usage:

  migoinfer [options] file.go [files.go...]

Options:

`
)

var (
	logPath   string
	showRaw   bool
	entryFunc string
	logFile   string
	logWriter = ioutil.Discard
)

func init() {
	flag.StringVar(&logPath, "log", "", "Specify analysis log file (use '-' for stderr)")
	flag.BoolVar(&showRaw, "raw", false, "Show raw unfiltered MiGo")
	flag.StringVar(&entryFunc, "entry", "", `Specify the function to view (format: (import/path).FuncName, empty means main.main)`)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, Usage)
		flag.PrintDefaults()
		os.Exit(0)
	}

	conf := build.FromFiles(flag.Args()...).Default()
	switch logPath {
	case "":
	case "-":
		logWriter = os.Stderr
		conf.WithBuildLog(logWriter, log.LstdFlags)
	default:
		f, err := os.Create(logPath)
		if err != nil {
			log.Fatalf("Cannot create log %s: %v", logPath, err)
		}
		defer f.Close()
		conf = conf.WithBuildLog(f, log.LstdFlags)
		logWriter = f
		logFile = f.Name()
	}
	info, err := conf.Build()
	if err != nil {
		log.Fatal("Build failed:", err)
	}
	inferer := migoinfer.New(info, logWriter)
	if logFile != "" {
		inferer.AddLogFiles(logFile)
	}
	if entryFunc != "" {
		inferer.SetEntryFunc(entryFunc)
	}
	inferer.SetOutput(os.Stdout)
	if showRaw {
		inferer.Raw = false
	}
	inferer.Analyse()
}
