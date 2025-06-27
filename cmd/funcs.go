package cmd

import (
	"fmt"
	"os"
)

func (us *CycloURL) AddURL(rawURL string) error {
	us.stats.total++

	domain, err := extractDomain(rawURL)
	if err != nil {
		us.stats.invalid++
		return err
	}

	us.mu.Lock()
	idx, exists := us.domainMap[domain]
	if !exists {
		idx = len(us.buckets)
		us.buckets = append(us.buckets, domainBucket{
			domain: domain,
			urls:   make([]string, 0, 8),
		})
		us.domainMap[domain] = idx
	}
	us.buckets[idx].urls = append(us.buckets[idx].urls, rawURL)
	us.stats.valid++
	us.mu.Unlock()

	return nil
}

func (us *CycloURL) InterleaveURLs() []string {
	us.mu.RLock()
	defer us.mu.RUnlock()

	if len(us.buckets) == 0 {
		return []string{}
	}

	result := make([]string, 0, us.stats.valid)
	lastUsed := make([]int, len(us.buckets))
	currentPos := make([]int, len(us.buckets))

	for i := range lastUsed {
		lastUsed[i] = -1
	}

	position := 0
	for {

		bestDomain := -1
		oldestUse := position

		for i := range us.buckets {
			if currentPos[i] < len(us.buckets[i].urls) {
				if lastUsed[i] < oldestUse {
					oldestUse = lastUsed[i]
					bestDomain = i
				}
			}
		}

		if bestDomain == -1 {
			break
		}

		result = append(result, us.buckets[bestDomain].urls[currentPos[bestDomain]])
		currentPos[bestDomain]++
		lastUsed[bestDomain] = position
		position++
	}

	return result
}

func (us *CycloURL) PrintStats(verbose bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	if verbose {
		fmt.Fprintf(os.Stderr, "URLs: %d total, %d valid, %d invalid\n",
			us.stats.total, us.stats.valid, us.stats.invalid)
		fmt.Fprintf(os.Stderr, "Domains: %d unique\n", len(us.buckets))

		if len(us.buckets) > 0 {
			fmt.Fprintf(os.Stderr, "Distribution:\n")
			for _, bucket := range us.buckets {
				fmt.Fprintf(os.Stderr, "  %s: %d URLs\n", bucket.domain, len(bucket.urls))
			}
		}
	} else if us.stats.valid > 0 {
		fmt.Fprintf(os.Stderr, "Processed %d URLs across %d domains\n", us.stats.valid, len(us.buckets))
	}
}
