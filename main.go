package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/zhhc99/bgen/internal/build"
	"github.com/zhhc99/bgen/internal/scaffold"
	"github.com/zhhc99/bgen/internal/server"
)

var Version = "dev"

func version() string {
	if Version != "dev" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return Version
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "init":
		err = scaffold.Run(".")
	case "build":
		err = build.Run(".")
	case "serve":
		err = server.Run(".")
	case "version":
		fmt.Println(version())
		return
	case "-h", "--help", "help":
		printUsage()
		return
	default:
		fmt.Fprintf(os.Stderr, "bgen: unknown command %q\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "bgen: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: bgen <command>")
	fmt.Fprintln(os.Stderr, "  init     initialize blog project scaffold")
	fmt.Fprintln(os.Stderr, "  build    build site to output/")
	fmt.Fprintln(os.Stderr, "  serve    start dev server with live reload")
	fmt.Fprintln(os.Stderr, "  version  print version")
}
