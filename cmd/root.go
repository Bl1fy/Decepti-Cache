package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Bl1fy/DeceptiCache/scanner"
	"github.com/spf13/cobra"
)

var (
	url            string
	urlFile        string
	headers        []string
	rateLimit      int
	onlyVulnerable bool
	requestRepeats int

	rootCmd = &cobra.Command{
		Use:   "decepticache",
		Short: "A tool for detecting Cache Deception vulnerabilities",
		Run: func(cmd *cobra.Command, args []string) {

			headerMap := convertHeaders(headers)

			if url != "" {
				fmt.Println("🔍 Scanning single URL...")
				scanner.ScanWCD(url, headerMap, onlyVulnerable, rateLimit, requestRepeats)
				return
			}

			if urlFile != "" {
				fmt.Println("Ohhh... This feature hasn't been implemented yet.")

				return
			}

			fmt.Println("❌ Please specify either --url or --urls")
			os.Exit(1)
		},
	}
)

func convertHeaders(headers []string) map[string]string {
	newHeaders := make(map[string]string)

	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			newHeaders[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return newHeaders
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Single URL to scan")
	rootCmd.Flags().StringVarP(&urlFile, "urls", "l", "", "File containing multiple URLs")
	rootCmd.Flags().StringArrayVarP(&headers, "header", "H", nil, "Custom HTTP headers")
	rootCmd.Flags().IntVarP(&rateLimit, "rate", "r", 10, "Maximum concurrent requests")
	rootCmd.Flags().BoolVarP(&onlyVulnerable, "only-vulnerable", "o", false, "Show only vulnerable results")
	rootCmd.Flags().IntVar(&requestRepeats, "request-repeats", 3, "How many times each payload should be repeated")
}
