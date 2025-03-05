package scanner

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type HTTPInfo struct {
	StatusCode int
	Headers    map[string]string
	Request    string
}

func GetUrlPath(RawUrl string) (string, error) {

	parsedUrl, err := url.Parse(RawUrl)

	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	return parsedUrl.Path, nil
}

func GetBaseURLAndPath(rawURL string) (string, string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	path := parsedURL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return baseURL, path, nil
}

func MakeRequest(url string, headers map[string]string) (*HTTPInfo, error) {
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("[!] Invalid URL: %s", url)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	var reqBuffer bytes.Buffer
	if err := req.Write(&reqBuffer); err != nil {
		return nil, fmt.Errorf("error writing request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	httpInfo := &HTTPInfo{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		Request:    reqBuffer.String(),
	}

	for key, values := range resp.Header {
		httpInfo.Headers[key] = strings.Join(values, ", ")
	}

	return httpInfo, nil
}

func CrawlLinks(pageURL string, headers map[string]string) ([]string, error) {
	var staticPaths []string
	seen := make(map[string]struct{})

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't receive response: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing document: %w", err)
	}

	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	extractPath := func(rawLink string) {
		if rawLink == "" || strings.HasPrefix(rawLink, "javascript:") {
			return
		}

		if strings.HasPrefix(rawLink, "/") {
			rawLink = baseURL + rawLink
		}

		parsed, err := url.Parse(rawLink)
		if err != nil || parsed.Host != parsedURL.Host {
			return
		}

		cleanPath := parsed.Path
		if _, found := seen[cleanPath]; !found {
			seen[cleanPath] = struct{}{}
			staticPaths = append(staticPaths, cleanPath)
		}
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			extractPath(link)
		}
	})

	doc.Find("link[href]").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			extractPath(link)
		}
	})

	doc.Find("script[src]").Each(func(i int, s *goquery.Selection) {
		if script, exists := s.Attr("src"); exists {
			extractPath(script)
		}
	})

	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		if img, exists := s.Attr("src"); exists {
			extractPath(img)
		}
	})

	return staticPaths, nil
}
