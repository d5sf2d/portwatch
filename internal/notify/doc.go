// Package notify defines the Notifier interface and built-in backends
// (log, webhook) used by portwatch to deliver change alerts.
//
// Use NewMulti to fan out a single Message to several backends at once.
package notify
