package globals

import (
	"regexp"
	"sync"

	"../response"
)

// package the globals away, I could have left them in the main file, but I didn't
// want to
var (
	EmailRE      *regexp.Regexp
	Header       []string
	Filters      []string
	Data         chan *response.HTTPResponse
	RequestMu    sync.Mutex // protects requestCount
	RequestCount int        // incremented on each request
)
