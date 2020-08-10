package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/k-yomo/fixtory"
	"golang.org/x/xerrors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	version   string
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name; default srcdir/fixtory_gen.go")
	pkgName   = flag.String("package", "", "package name; default same package as the type")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of fixtory %s:\n", version)
	fmt.Fprintf(os.Stderr, "\tfixtory [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("fixtory: ")
	flag.Usage = Usage

	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	types := strings.Split(*typeNames, ",")

	args := flag.Args()
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	outputPath := *output
	if outputPath == "" {
		outputPath = filepath.Join(targetDir, "fixtory_gen.go")
	}

	newWriter := func() (io.Writer, func() error, error) {
		writer, err := os.Create(outputPath)
		if err != nil {
			return nil, nil, xerrors.Errorf("create output file: %w", err)
		}
		return writer, func() error { return writer.Close() }, nil
	}
	if err := fixtory.Generate(targetDir, types, pkgName, newWriter); err != nil {
		color.Red("%+v", err)
		os.Exit(1)
	}
}
