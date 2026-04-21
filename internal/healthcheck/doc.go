// Package healthcheck probes open TCP ports to confirm they are accepting
// connections. It is used by the pipeline to annotate scan results with
// liveness information before alerts are dispatched.
package healthcheck
