package shared

import (
	"encoding/json"
	"html"
	"strings"
	"time"
)

const commentsKey = "comments"

// MaxComments caps how many letters we keep, oldest dropped first, so the
// single JSON blob never grows unbounded on a free Redis tier.
const MaxComments = 500

type Comment struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

// LoadComments returns all stored comments, oldest first. If Redis is
// unreachable or not yet configured, it returns an empty list rather than
// an error - the story itself should still be readable either way.
func LoadComments() ([]Comment, error) {
	raw, ok, err := Get(commentsKey)
	if err != nil || !ok {
		return []Comment{}, nil
	}
	var list []Comment
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return []Comment{}, nil
	}
	return list, nil
}

func saveComments(list []Comment) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return Set(commentsKey, string(b))
}

// sanitize trims whitespace, hard-caps length, and HTML-escapes free text so
// nothing a reader submits can inject markup into the page.
func sanitize(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	// Collapse excessive internal whitespace/newlines a little.
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > maxLen {
		s = s[:maxLen]
	}
	return html.EscapeString(s)
}

// AddComment validates, sanitizes, and appends a new reader letter.
func AddComment(name, text string) (Comment, error) {
	name = sanitize(name, 60)
	text = sanitize(text, 2000)

	if name == "" {
		name = "A Reader"
	}

	list, err := LoadComments()
	if err != nil {
		return Comment{}, err
	}

	c := Comment{
		ID:        time.Now().UTC().Format("20060102150405.000000"),
		Name:      name,
		Text:      text,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	list = append(list, c)
	if len(list) > MaxComments {
		list = list[len(list)-MaxComments:]
	}
	if err := saveComments(list); err != nil {
		return Comment{}, err
	}
	return c, nil
}

// DeleteComment removes a comment by ID and reports whether it was found.
func DeleteComment(id string) (bool, error) {
	list, err := LoadComments()
	if err != nil {
		return false, err
	}
	out := make([]Comment, 0, len(list))
	found := false
	for _, c := range list {
		if c.ID == id {
			found = true
			continue
		}
		out = append(out, c)
	}
	if !found {
		return false, nil
	}
	return true, saveComments(out)
}
