// Package portclassify categorises ports into well-known service tiers
// (system, registered, dynamic) and marks them as privileged or unprivileged.
package portclassify

// Tier represents the IANA port range classification.
type Tier string

const (
	// TierSystem covers ports 0–1023 (well-known / privileged).
	TierSystem Tier = "system"
	// TierRegistered covers ports 1024–49151.
	TierRegistered Tier = "registered"
	// TierDynamic covers ports 49152–65535 (ephemeral).
	TierDynamic Tier = "dynamic"
)

// Result holds the classification of a single port.
type Result struct {
	Port       int
	Tier       Tier
	Privileged bool // true when port < 1024
}

// Classifier classifies ports into tiers.
type Classifier struct{}

// New returns a ready-to-use Classifier.
func New() *Classifier {
	return &Classifier{}
}

// Classify returns the classification Result for the given port number.
// Invalid port numbers (< 0 or > 65535) return an error.
func (c *Classifier) Classify(port int) (Result, error) {
	if port < 0 || port > 65535 {
		return Result{}, fmt.Errorf("portclassify: port %d out of range [0, 65535]", port)
	}

	r := Result{
		Port:       port,
		Privileged: port < 1024,
	}

	switch {
	case port <= 1023:
		r.Tier = TierSystem
	case port <= 49151:
		r.Tier = TierRegistered
	default:
		r.Tier = TierDynamic
	}

	return r, nil
}

// ClassifyAll classifies a slice of port numbers and returns a slice of
// Results in the same order. Ports that are out of range are skipped.
func (c *Classifier) ClassifyAll(ports []int) []Result {
	out := make([]Result, 0, len(ports))
	for _, p := range ports {
		r, err := c.Classify(p)
		if err != nil {
			continue
		}
		out = append(out, r)
	}
	return out
}
