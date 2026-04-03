package api

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	DebugNone    = 0
	DebugAPI     = 1
	DebugVerbose = 2
)

var debugLevel int

func SetDebugLevel(level int) {
	debugLevel = level
}

func debugLogRequest(req *http.Request) {
	if debugLevel < DebugAPI {
		return
	}
	fmt.Fprintf(os.Stderr, ">> %s %s\n", req.Method, req.URL)
	if debugLevel >= DebugVerbose {
		for k, v := range req.Header {
			fmt.Fprintf(os.Stderr, ">> %s: %s\n", k, v)
		}
	}
}

func debugLogResponse(resp *http.Response, elapsed time.Duration) {
	if debugLevel < DebugAPI {
		return
	}
	fmt.Fprintf(os.Stderr, "<< %d %s (%s)\n", resp.StatusCode, resp.Status, elapsed)
}
