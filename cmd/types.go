package cmd

import "sync"

type Config struct {
	inputFile  string
	outputFile string
	bufferSize int
	concurrent int
	showStats  bool
	validate   bool
	verbose    bool
	stdin      bool
}

type CycloURL struct {
	buckets   []domainBucket
	domainMap map[string]int
	stats     struct {
		total,
		valid,
		invalid int
	}
	mu sync.RWMutex
}

type domainBucket struct {
	domain string
	urls   []string
}
