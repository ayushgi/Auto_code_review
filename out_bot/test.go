package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	pollSleep = time.Sleep
	pollDone  = func() {}
)

type Server struct {
	version string
	url     string
	period  time.Duration
	tagged  bool
}

func NewServer(version, url string, period time.Duration) *Server {
	s := &Server{
		version: version,
		url:     url,
		period:  period,
	}
	go s.poll()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.tagged {
		fmt.Fprintf(w, "YES! Version %s is tagged.\n", s.version)
		return
	}
	fmt.Fprintf(w, "No. Version %s is not tagged.\n", s.version)
}

func (s *Server) poll() {
	for !isTagged(s.url) {
		pollSleep(s.period)
	}
	s.tagged = true
	pollDone()
}

func isTagged(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
} 