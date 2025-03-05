package scanner

import (
	"fmt"
	"path"
	"strings"
	"sync"
	"time"
)

var cacheHeaders = []string{
	"X-Cache", "CF-Cache-Status", "Age", "CDN-Cache-Control",
	"Server-Timing", "Cache-Control", "Pragma", "Surrogate-Control",
	"Vary", "Expires",
}

var onlyVulnerable = false

var cachedValue = []string{"hit", "cached", "store"}

func isCached(headers map[string]string) bool {
	for _, header := range cacheHeaders {
		if value, exists := headers[header]; exists {
			lowerValue := strings.ToLower(value)

			if strings.Contains(lowerValue, "no-cache") || strings.Contains(lowerValue, "no-store") || strings.Contains(lowerValue, "private") {
				return false
			}

			for _, cached := range cachedValue {
				if strings.Contains(lowerValue, cached) {
					return true
				}
			}
		}
	}
	return false
}

func getCacheHeaders(headers map[string]string) string {
	var returnHeader = []string{}

	for _, headerName := range cacheHeaders {
		if value, exists := headers[headerName]; exists {
			returnHeader = append(returnHeader, fmt.Sprintf("%s: %s", headerName, value))
		}
	}

	if len(returnHeader) > 0 {
		return strings.Join(returnHeader, " | ")
	}

	return "None"
}

func constructTestURL(baseURL, payload string) string {
	baseURL = strings.TrimRight(baseURL, "/")

	if strings.HasPrefix(payload, "/") {
		return baseURL + payload
	}
	return baseURL + payload
}

func testPayloads(url string, headers map[string]string, payloads []string, maxConcurrency int, requestRepeats int, onlyVulnerable bool) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)

	for _, payload := range payloads {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			testURL := constructTestURL(url, p)
			var isPayloadCached bool
			var testRequest *HTTPInfo
			var lastErr error
			var err error

			for i := 0; i < requestRepeats; i++ {
				testRequest, err = MakeRequest(testURL, headers)
				if err != nil {
					lastErr = err
					time.Sleep(time.Duration(25*(1<<i)) * time.Millisecond)
					continue
				}

				if testRequest == nil {
					lastErr = fmt.Errorf("no response received")
					continue
				}

				isPayloadCached = isCached(testRequest.Headers)
				if isPayloadCached {
					break
				}
			}

			if testRequest == nil {
				fmt.Printf("[ERROR] %s -> %v\n", testURL, lastErr)
				return
			}

			if isPayloadCached {
				fmt.Printf("[🔥 CACHED] %s -> %s | Status: %d | Cache-Header: %s\n",
					url, testURL, testRequest.StatusCode, getCacheHeaders(testRequest.Headers))
			} else if !onlyVulnerable {
				fmt.Printf("[✅ SAFE] %s -> %s | Status: %d | Cache-Header: %s\n",
					url, testURL, testRequest.StatusCode, getCacheHeaders(testRequest.Headers))
			}
		}(payload)
	}

	wg.Wait()
}

func testEXT(url string, headers map[string]string, maxConcurrency int, requestRepeats int, onlyVulnerable bool) {
	var payloads []string
	for _, extension := range Extensions {
		for _, payload := range ExtPayloads {
			payloads = append(payloads, payload+extension)
		}
	}
	testPayloads(url, headers, payloads, maxConcurrency, requestRepeats, onlyVulnerable)
}

func testDelimeters(url string, headers map[string]string, maxConcurrency int, requestRepeats int, onlyVulnerable bool) {
	var payloads []string

	for _, file := range CommonFiles {
		for _, delim := range Delimiters {
			payloads = append(payloads, delim+file)
			payloads = append(payloads, delim+delim+file)
			payloads = append(payloads, "/.."+delim+file)
			payloads = append(payloads, delim+"/../"+file)
		}

		payloads = append(payloads, "?"+file)
		payloads = append(payloads, ";"+file)
		payloads = append(payloads, "%3b"+file)
		payloads = append(payloads, "?_="+file)
		payloads = append(payloads, "?version="+file)
	}

	testPayloads(url, headers, payloads, maxConcurrency, requestRepeats, onlyVulnerable)
}

func normalizeStaticPaths(crawledLinks []string) []string {
	var staticPaths []string

	for _, link := range crawledLinks {
		if !strings.HasPrefix(link, "/") {
			continue
		}

		dirPath := link
		if path.Ext(link) != "" {
			dirPath = path.Dir(link)
		}

		for dirPath != "/" && dirPath != "." && dirPath != "" {
			if dirPath == "/" || dirPath == "." || dirPath == "" {
				break
			}

			staticPaths = append(staticPaths, dirPath)
			newPath := path.Clean(path.Dir(dirPath))

			if newPath == dirPath {
				break
			}

			dirPath = newPath
		}
	}

	return staticPaths
}

func testStaticPath(url string, headers map[string]string, maxConcurrency int, requestRepeats int, onlyVulnerable bool) {
	crawled, err := CrawlLinks(url, headers)

	if err != nil {
		crawled = append(crawled, "")
	}

	staticPaths := normalizeStaticPaths(crawled)

	host, relevantPath, err := GetBaseURLAndPath(url)

	relevantPath = relevantPath[1:]

	if err != nil {
		relevantPath = "/"
	}

	var payloads []string
	for _, path := range staticPaths {
		for _, nsp := range normalizeStaticPathsPayloads {
			payloads = append(payloads, path+nsp+relevantPath)
		}
	}

	testPayloads(host, headers, payloads, maxConcurrency, requestRepeats, onlyVulnerable)
}

func ScanWCD(url string, headers map[string]string, show_onlyVulnerable bool, maxConcurrency int, requestRepeats int) {
	onlyVulnerable = show_onlyVulnerable

	originReqInfo, err := MakeRequest(url, headers)
	if err != nil {
		fmt.Printf("[ERROR] %s - %v\n", url, err)
		return
	}

	if originReqInfo == nil {
		fmt.Printf("[ERROR] No response received from %s\n", url)
		return
	}

	isOriginCached := isCached(originReqInfo.Headers)
	fmt.Printf("🔍 Testing: %s | Status: %d | Cached: %v\n", url, originReqInfo.StatusCode, isOriginCached)

	if isOriginCached && onlyVulnerable {
		fmt.Printf("[SKIPPED] %s - Already cached, skipping further tests.\n", url)
		return
	}

	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		testEXT(url, headers, maxConcurrency, requestRepeats, onlyVulnerable)
	}()
	go func() {
		defer wg.Done()
		testDelimeters(url, headers, maxConcurrency, requestRepeats, onlyVulnerable)
	}()
	go func() {
		defer wg.Done()
		testStaticPath(url, headers, maxConcurrency, requestRepeats, onlyVulnerable)
	}()

	wg.Wait()
}
