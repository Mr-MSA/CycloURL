package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	version           = "1.0.0"
	defaultBufferSize = 64 * 1024
	maxConcurrency    = 16
	batchSize         = 1000
)

var (
	urlPool = sync.Pool{New: func() interface{} { return make([]string, 0, 16) }}
	sbPool  = sync.Pool{New: func() interface{} { return &strings.Builder{} }}
)

func NewCycloURL(estimatedDomains int) *CycloURL {

	if estimatedDomains <= 0 {
		estimatedDomains = 64
	}

	return &CycloURL{
		buckets:   make([]domainBucket, 0, estimatedDomains),
		domainMap: make(map[string]int, estimatedDomains),
	}
}

func Execute() {
	cfg := parseConfig()

	input, inputCleanup, err := openInput(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer inputCleanup()

	output, outputCleanup, err := createOutput(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer outputCleanup()

	estimatedURLs := estimateFileSize(cfg.inputFile)
	sorter := NewCycloURL(estimatedURLs / 10)

	if err := processURLs(input, sorter, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading URLs: %v\n", err)
		os.Exit(1)
	}

	if sorter.stats.valid == 0 {
		fmt.Fprintln(os.Stderr, "No valid URLs found")
		os.Exit(1)
	}

	sortedURLs := sorter.InterleaveURLs()

	if err := writeURLs(sortedURLs, output, cfg.bufferSize); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing URLs: %v\n", err)
		os.Exit(1)
	}

	if cfg.showStats || cfg.verbose {
		sorter.PrintStats(cfg.verbose)
	}
}
