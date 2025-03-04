package handlers

import (
	"fmt"
	"time"
)

// TimeAgo returns a human-readable string representing the time elapsed since the given time.
func TimeAgo(createdAt time.Time) string {
	now := time.Now()
	diff := now.Sub(createdAt)

	switch {
	case diff < time.Second:
		return "just now"
	case diff < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	default:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}