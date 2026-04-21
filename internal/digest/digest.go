// Package digest produces a stable hash of a port snapshot so callers
// can cheaply detect whether anything changed between two scans.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/state"
)

// Digest holds the hex-encoded SHA-256 hash of a snapshot.
type Digest struct {
	Hash string
}

// Equal returns true when two Digests represent identical snapshots.
func (d Digest) Equal(other Digest) bool {
	return d.Hash == other.Hash
}

// String implements fmt.Stringer.
func (d Digest) String() string { return d.Hash }

// Of computes a deterministic Digest from a snapshot.
// Ports are sorted before hashing so insertion order does not matter.
func Of(snap state.Snapshot) Digest {
	ports := make([]int, 0, len(snap.Ports))
	for _, p := range snap.Ports {
		ports = append(ports, p)
	}
	sort.Ints(ports)

	h := sha256.New()
	for _, p := range ports {
		fmt.Fprintf(h, "%d\n", p)
	}
	return Digest{Hash: hex.EncodeToString(h.Sum(nil))}
}

// Changed returns true when the two snapshots produce different Digests.
func Changed(prev, curr state.Snapshot) bool {
	return !Of(prev).Equal(Of(curr))
}
