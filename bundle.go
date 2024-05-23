package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/evanw/esbuild/pkg/api"
)

func Push[T any](s *[]T, v T) {
	if *s == nil {
		*s = make([]T, 0, 1)
	}
	if len(*s) == cap(*s) {
		*s = append(make([]T, 0, cap(*s)*2), *s...)
	}
	*s = append(*s, v)
}

func main() {
	var minify bool
	var inFile, outFile []string
	startTime := time.Now()
	flag.Func("in", "Input file to bundle and optionally minify", func(in string) error {
		if len(inFile) > len(outFile) {
			return fmt.Errorf("Error: expecting -in option")
		}
		Push(&inFile, in)
		return nil
	})
	flag.Func("out", "Bundled and optionally minified output goes into this file", func(out string) error {
		if len(outFile) > len(inFile) {
			return fmt.Errorf("Error: expecting -out option")
		}
		Push(&outFile, out)
		return nil
	})
	flag.BoolVar(&minify, "minify", false, "Minify output (default: false)")
	flag.Parse()

	if len(inFile) == 0 || len(outFile) == 0 || len(inFile) != len(outFile) {
		fmt.Fprint(os.Stderr, "Error: must specify one or more matching -in and -out options\n")
		fmt.Fprintf(os.Stderr, "Usage: %s options\nAvailable options are:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	for i, _ := range inFile {
		result := api.Build(api.BuildOptions{
			EntryPoints:       []string{inFile[i]},
			Outfile:           outFile[i],
			Bundle:            true,
			MinifyWhitespace:  minify,
			MinifyIdentifiers: minify,
			MinifySyntax:      minify,
			Write:             true,
			TreeShaking:       api.TreeShakingTrue,
			LogLevel:          api.LogLevelWarning,
		})

		if len(result.Errors) > 0 {
			os.Exit(1)
		}
	}
	fmt.Fprintf(os.Stderr, "Bundled in: %s\n", time.Now().Sub(startTime))
}
