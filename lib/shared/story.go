package shared

import (
	_ "embed"
	"strings"
)

// defaultStory is the manuscript as originally supplied, embedded into the
// binary at build time. It is served as-is until the author overwrites it
// through the admin panel, at which point the Redis-stored version wins.
//
//go:embed story_data/story.txt
var defaultStory string

const storyKey = "story_content"
const storyTitleKey = "story_title"

const defaultTitle = "The Man Who Threw Newspapers"

// GetStory returns the current title and body, preferring whatever the
// author has saved in Redis and falling back to the original manuscript.
// If Redis is unreachable or not yet configured, it still returns the
// embedded default rather than failing the whole page - a missing/misconfigured
// database should never be the reason a reader can't read the story.
func GetStory() (title string, body string, err error) {
	title = defaultTitle
	body = defaultStory

	if t, ok, e := Get(storyTitleKey); e == nil && ok && strings.TrimSpace(t) != "" {
		title = t
	}

	if b, ok, e := Get(storyKey); e == nil && ok && strings.TrimSpace(b) != "" {
		body = b
	}

	return title, body, nil
}

// SaveStory persists a new title/body pair, replacing the manuscript.
func SaveStory(title, body string) error {
	if strings.TrimSpace(title) != "" {
		if err := Set(storyTitleKey, title); err != nil {
			return err
		}
	}
	return Set(storyKey, body)
}
