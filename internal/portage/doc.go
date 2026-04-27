// Package portage tracks the continuous open-duration of ports on monitored
// hosts. It records when each port was first seen open and exposes its age,
// making it easy to detect long-lived or stale services that should be
// reviewed.
package portage
