package scanner

var Extensions = []string{
	"css", "js", "jpg", "png", "gif", "ico", "woff", "woff2", "ttf", "eot", "svg", "otf",
	"mp3", "mp4", "webm", "ogg", "wav", "flac", "avi", "mov", "mkv",
	"zip", "rar", "tar", "gz", "bz2", "xz", "7z", "iso",
	"doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf", "rtf",
	"class", "jar", "war", "ear", "exe", "dll", "so", "deb", "rpm",
	"json", "xml", "csv", "txt", "log", "conf", "yaml", "yml", "ini",
}

var ExtPayloads = []string{
	"/.", "/%2e", "%2f.", "%2f%2e", "/test.", "/тест.", "/test%2e", "/тест%2e", "/测试.", "/测试%2e",
	"%2ftest.", "%2fтест.", "%2f测试.", "/%252e", "/%252f",
}

var Delimiters = []string{
	"/", "/.", "/..", ";", "?", "#", "%2e", "%2f", "%3b", "%23",
	"%3f", "%c0%af", "%e3%80%82", "%ef%bc%8f", "%00",
}
var normalizeStaticPathsPayloads = []string{
	"/../", "../", "/../../", "../../", "/../../../", "../../../",
	"/..%2f", "..%2f", "/../..%2f", "../..%2f", "/../../..%2f", "../../..%2f",
	"/..%2f", "..%2f", "/../..%2f", "../..%2f", "/../../..%2f", "../../..%2f",
}

var StaticPathPayloads = []string{
	"/static", "/assets", "/cdn-cgi", "/media", "/resources",
}

var CommonFiles = []string{
	"robots.txt", "index.html", "index.php", "sitemap.xml", "favicon.ico",
	"style.css", "script.js",
}
