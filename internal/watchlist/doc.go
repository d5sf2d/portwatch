// Package watchlist provides named port group management for portwatch.
//
// A Watchlist organises ports into logical groups (e.g. "web", "database").
// Groups can be validated, merged into a flat port list for scanning, and
// queried to determine which groups own a specific port.
package watchlist
