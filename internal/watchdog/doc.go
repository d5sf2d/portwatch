// Package watchdog detects "flapping" ports — ports whose open/closed state
// changes more than a configurable number of times within a rolling time
// window.  When a breach is detected the caller receives a Breach value that
// can be forwarded to any Notifier.
package watchdog
