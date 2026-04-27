// Package portpriority provides priority level classification for ports.
// Levels range from Low to Critical and are used by the pipeline to
// order alerts and filter noise based on operator-defined thresholds.
//
// Priority levels are ordered as follows (lowest to highest):
//
//	Low < Medium < High < Critical
//
// Typical usage involves assigning a Level to each observed port and
// comparing it against a configured minimum threshold to decide whether
// an alert should be emitted or suppressed.
package portpriority
