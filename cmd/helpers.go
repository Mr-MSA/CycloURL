package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
)

func processURLs(input io.Reader, sorter *CycloURL, cfg *Config) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, cfg.bufferSize), cfg.bufferSize*2)

	if cfg.concurrent == 1 {
		for scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" && !strings.HasPrefix(line, "#") {
				if err := sorter.AddURL(line); err != nil && cfg.validate {
					if cfg.verbose {
						fmt.Fprintf(os.Stderr, "Invalid URL: %v\n", err)
					}
				}
			}
		}
	} else {
		semaphore := make(chan struct{}, cfg.concurrent)
		var wg sync.WaitGroup

		for scanner.Scan() {
			if line := strings.TrimSpace(scanner.Text()); line != "" && !strings.HasPrefix(line, "#") {
				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					semaphore <- struct{}{}
					defer func() { <-semaphore }()

					if err := sorter.AddURL(url); err != nil && cfg.validate && cfg.verbose {
						fmt.Fprintf(os.Stderr, "Invalid URL: %v\n", err)
					}
				}(line)
			}
		}
		wg.Wait()
	}

	return scanner.Err()
}

func estimateFileSize(filename string) int {
	if filename == "-" {
		return 1000
	}
	if stat, err := os.Stat(filename); err == nil {
		return max(100, int(stat.Size())/50)
	}
	return 1000
}

func extractDomain(rawURL string) (string, error) {
	switch {
	case strings.HasPrefix(rawURL, "//"):
		rawURL = "https:" + rawURL
	case !strings.Contains(rawURL, "://"):
		rawURL = "https://" + rawURL
	}

	if parsed, err := url.Parse(rawURL); err != nil {
		return "", err
	} else if parsed.Host == "" {
		return "unknown", nil
	} else {
		return parsed.Host, nil
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func hasStdinInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) == 0
}

// io
func writeURLs(urls []string, output io.Writer, bufferSize int) error {
	writer := bufio.NewWriterSize(output, bufferSize)
	defer writer.Flush()

	sb := sbPool.Get().(*strings.Builder)
	defer func() { sb.Reset(); sbPool.Put(sb) }()

	for i := 0; i < len(urls); i += batchSize {
		end := min(i+batchSize, len(urls))
		sb.Reset()

		for j := i; j < end; j++ {
			sb.WriteString(urls[j])
			sb.WriteByte('\n')
		}

		if _, err := writer.WriteString(sb.String()); err != nil {
			return err
		}
	}
	return nil
}

func createOutput(cfg *Config) (io.Writer, func(), error) {
	if cfg.outputFile == "" {
		return os.Stdout, func() {}, nil
	}

	file, err := os.Create(cfg.outputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("creating output file: %w", err)
	}

	// Create MultiWriter to write to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, file)
	cleanup := func() { file.Close() }

	return multiWriter, cleanup, nil
}

func openInput(cfg *Config) (io.Reader, func(), error) {
	if cfg.stdin {
		return os.Stdin, func() {}, nil
	}

	file, err := os.Open(cfg.inputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("opening input file: %w", err)
	}

	return file, func() { file.Close() }, nil
}
