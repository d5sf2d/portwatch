// Package fingerprint identifies services running on open ports by
// reading banners or matching known response patterns.
package fingerprint

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Service holds the result of a fingerprint probe.
type Service struct {
	Port    int
	Banner  string
	Guess   string
}

// Fingerprinter probes ports and returns service guesses.
type Fingerprinter struct {
	Timeout time.Duration
}

// New returns a Fingerprinter with the given timeout.
func New(timeout time.Duration) *Fingerprinter {
	return &Fingerprinter{Timeout: timeout}
}

// Probe connects to host:port, reads a banner if available, and
// returns a Service with a best-effort service guess.
func (f *Fingerprinter) Probe(host string, port int) Service {
	svc := Service{Port: port}
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, f.Timeout)
	if err != nil {
		svc.Guess = "unknown"
		return svc
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(f.Timeout))
	buf := make([]byte, 256)
	n, _ := conn.Read(buf)
	if n > 0 {
		svc.Banner = strings.TrimSpace(string(buf[:n]))
	}
	svc.Guess = guessService(port, svc.Banner)
	return svc
}

// ProbeAll probes a list of ports on the given host.
func (f *Fingerprinter) ProbeAll(host string, ports []int) []Service {
	results := make([]Service, 0, len(ports))
	for _, p := range ports {
		results = append(results, f.Probe(host, p))
	}
	return results
}

func guessService(port int, banner string) string {
	b := strings.ToLower(banner)
	switch {
	case strings.Contains(b, "ssh"):
		return "ssh"
	case strings.Contains(b, "http") || strings.Contains(b, "html"):
		return "http"
	case strings.Contains(b, "ftp"):
		return "ftp"
	case strings.Contains(b, "smtp") || strings.Contains(b, "220"):
		return "smtp"
	}
	switch port {
	case 22:
		return "ssh"
	case 80, 8080, 443:
		return "http"
	case 21:
		return "ftp"
	case 25, 587:
		return "smtp"
	case 3306:
		return "mysql"
	case 5432:
		return "postgres"
	case 6379:
		return "redis"
	}
	return "unknown"
}
