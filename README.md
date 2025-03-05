# DeceptiCache
**DeceptiCache** is a tool designed to help detect **Web Cache Deception** vulnerabilities by automating payload execution.


# Installation

```
go install github.com/Bl1fy/Decepti-Cache@latest
```
OR
```
git clone https://github.com/Bl1fy/Decepti-Cache.git
cd Decepti-Cache
go build -o decepticache
```

# Usage
```
A tool for detecting Cache Deception vulnerabilities

Usage:
  decepticache [flags]

Flags:
  -H, --header stringArray    Custom HTTP headers
  -h, --help                  help for decepticache
  -o, --only-vulnerable       Show only vulnerable results
  -r, --rate int              Maximum concurrent requests (default 10)
      --request-repeats int   How many times each payload should be repeated (default 3)
  -u, --url string            Single URL to scan
  -l, --urls string           File containing multiple URLs
```

### Example Usage
```
$ decepticache --url "https://web-security-academy.net/my-account" -H "Cookie: session=43tqN0ZHf6KB93jKnJrqJ36kWXOgeHsK" --request-repeats 7 --rate 75 --only-vulnerable
```
