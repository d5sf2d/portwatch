package portschedule

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// WindowConfig is the JSON-serialisable form of a Window.
type WindowConfig struct {
	Start string `json:"start"` // e.g. "09:00"
	End   string `json:"end"`   // e.g. "17:00"
}

// Config is the JSON-serialisable form of a Schedule.
type Config struct {
	Windows  []WindowConfig `json:"windows"`
	Weekdays []string       `json:"weekdays"` // e.g. ["Monday","Tuesday"]
}

var weekdayNames = map[string]time.Weekday{
	"Sunday": time.Sunday, "Monday": time.Monday, "Tuesday": time.Tuesday,
	"Wednesday": time.Wednesday, "Thursday": time.Thursday,
	"Friday": time.Friday, "Saturday": time.Saturday,
}

// LoadFile reads a Schedule from a JSON file.
func LoadFile(path string) (*Schedule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("portschedule: read %s: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("portschedule: parse %s: %w", path, err)
	}
	return FromConfig(cfg)
}

// FromConfig converts a Config into a Schedule.
func FromConfig(cfg Config) (*Schedule, error) {
	var windows []Window
	for _, wc := range cfg.Windows {
		start, err := parseClock(wc.Start)
		if err != nil {
			return nil, err
		}
		end, err := parseClock(wc.End)
		if err != nil {
			return nil, err
		}
		windows = append(windows, Window{Start: start, End: end})
	}
	var days []time.Weekday
	for _, name := range cfg.Weekdays {
		d, ok := weekdayNames[name]
		if !ok {
			return nil, fmt.Errorf("portschedule: unknown weekday %q", name)
		}
		days = append(days, d)
	}
	return New(windows, days)
}

func parseClock(s string) (time.Duration, error) {
	var h, m int
	if _, err := fmt.Sscanf(s, "%d:%d", &h, &m); err != nil {
		return 0, fmt.Errorf("portschedule: invalid time %q", s)
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("portschedule: time out of range %q", s)
	}
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}
