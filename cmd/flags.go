package cmd

import (
	"fmt"
	"os"
	"runtime"

	flag "github.com/spf13/pflag"
)

func parseConfig() *Config {
	cfg := &Config{}
	var showVersion bool

	flag.StringVarP(&cfg.inputFile, "input", "i", "", "Input file containing URLs")
	flag.StringVarP(&cfg.outputFile, "output", "o", "", "Output file for sorted URLs")
	flag.BoolVarP(&cfg.showStats, "stats", "s", false, "Show basic processing statistics")
	flag.BoolVar(&cfg.verbose, "verbose", false, "Show detailed statistics and progress")
	flag.BoolVarP(&cfg.validate, "validate", "v", false, "Validate URLs and skip invalid ones")
	flag.IntVarP(&cfg.concurrent, "concurrent", "c", 1, "Max concurrent processing (1-16)")
	flag.IntVar(&cfg.bufferSize, "buffer", defaultBufferSize, "I/O buffer size in bytes")
	flag.BoolVar(&showVersion, "version", false, "Show version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "# CycloURL v%s\n", version)
		fmt.Fprintf(os.Stderr, "CLI Tool for distributing URLs across domains using smart round-robin method\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -i urls.txt -o sorted.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s < urls.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  cat urls.txt | %s -s -c 4\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --input large.txt --concurrent 8 --verbose\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFor more information, visit: https://github.com/Mr-MSA/CycloURL\n")
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("CycloURL %s\n", version)
		fmt.Printf("Go: %s, Platform: %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Max concurrency: %d, Default buffer: %d bytes\n", maxConcurrency, defaultBufferSize)
		os.Exit(0)
	}

	if cfg.inputFile == "" {

		if !hasStdinInput() {
			fmt.Fprintln(os.Stderr, "No input file specified. Use -i or pipe input to stdin")
			flag.Usage()
			os.Exit(1)
		} else {
			cfg.inputFile = "-"
		}

	} else {
		if _, err := os.Stat(cfg.inputFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Input file '%s' does not exist\n", cfg.inputFile)
			os.Exit(1)
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing input file '%s': %v\n", cfg.inputFile, err)
			os.Exit(1)
		}
	}

	cfg.stdin = cfg.inputFile == "-"

	if cfg.concurrent < 1 || cfg.concurrent > maxConcurrency {
		cfg.concurrent = min(4, runtime.NumCPU())
	}
	if cfg.bufferSize < 1024 {
		cfg.bufferSize = 1024
	}

	return cfg
}
