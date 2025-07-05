# CycloURL
A CLI tool that distributes URLs across domains using round-robin scheduling to help prevent rate limiting when making sequential requests

## Installation
### Direct Install
```
go install github.com/Mr-MSA/cyclourl@latest
```
### From Source
```
git clone https://github.com/Mr-MSA/cyclourl
cd cyclourl
go build . -o cyclourl
```
## Usage
### Basic
```
# Sort URLs from file
./cyclourl -i urls.txt -o sorted.txt

# Read from stdin, output to stdout
cat urls.txt | ./cyclourl

# Show statistics
./cyclourl -i urls.txt -s
```
### Optimized
```
# High-performance processing with concurrency
./cyclourl -i large_file.txt -c 8 --buffer 131072

# Pipeline processing with validation  
curl -s https://example.com/urls.txt | ./cyclourl -i - -v -c 4

# Memory-efficient processing of huge files
./cyclourl -i million_urls.txt -o distributed.txt -c 8 -s
```
### Flags
+ Help output: `$ ./cyclourl --help`
```
      --buffer int       I/O buffer size in bytes (default 65536)
  -c, --concurrent int   Max concurrent processing (1-16) (default 1)
  -i, --input string     Input file containing URLs
  -o, --output string    Output file for sorted URLs
  -s, --stats            Show basic processing statistics
  -v, --validate         Validate URLs and skip invalid ones
      --verbose          Show detailed statistics and progress
      --version          Show version information
```

## License
MIT License - see [LICENSE](https://github.com/Mr-MSA/CycloURL/blob/main/LICENSE) file for details
