package shared

import (
	"strings"
	"testing"
)

func TestGetStoryFallsBackToEmbeddedManuscriptWithoutRedis(t *testing.T) {
	title, body, err := GetStory()
	if err != nil {
		t.Fatalf("expected no error even without Redis configured, got: %v", err)
	}
	if title != defaultTitle {
		t.Errorf("expected default title %q, got %q", defaultTitle, title)
	}
	if !strings.Contains(body, "Karthick") {
		t.Errorf("expected embedded manuscript to be returned, got body of length %d", len(body))
	}
}

func TestLoadCommentsWithoutRedisReturnsEmptyNotError(t *testing.T) {
	list, err := LoadComments()
	if err != nil {
		t.Fatalf("expected no error even without Redis configured, got: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}

func TestSanitizeEscapesHTML(t *testing.T) {
	got := sanitize("<script>alert(1)</script>", 100)
	if strings.Contains(got, "<script>") {
		t.Errorf("expected HTML to be escaped, got %q", got)
	}
}

func TestSanitizeTruncatesLongInput(t *testing.T) {
	long := strings.Repeat("a", 200)
	got := sanitize(long, 50)
	if len(got) > 50 {
		t.Errorf("expected truncation to 50 chars, got %d", len(got))
	}
}
