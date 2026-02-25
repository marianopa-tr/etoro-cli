package output

import (
	"encoding/json"
	"testing"
)

func TestPrintFeedPostsTable(t *testing.T) {
	posts := []FeedPostRow{
		{ID: "p1", Author: "trader1", Message: "Bullish on AAPL!", Created: "2025-01-15"},
		{ID: "p2", Author: "trader2", Message: "This is a very long message that should be truncated because it exceeds the maximum length allowed in the table display which is eighty characters", Created: "2025-01-14"},
	}

	out := captureStdout(t, func() {
		PrintFeedPosts(posts, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintFeedPostsJSON(t *testing.T) {
	posts := []FeedPostRow{{ID: "p1", Author: "trader1", Message: "Hello"}}

	out := captureStdout(t, func() {
		PrintFeedPosts(posts, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintFeedPostsEmpty(t *testing.T) {
	captureStdout(t, func() {
		PrintFeedPosts(nil, Table)
	})
}

func TestPrintPostCreatedTable(t *testing.T) {
	PrintPostCreated(Table)
}

func TestPrintPostCreatedJSON(t *testing.T) {
	out := captureStdout(t, func() {
		PrintPostCreated(JSON)
	})

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["status"] != "created" {
		t.Errorf("status = %v", result["status"])
	}
}

func TestPrintFeedPostSingle(t *testing.T) {
	post := FeedPostRow{ID: "p1", Author: "trader1", Message: "Bullish!", Created: "2025-01-15"}

	out := captureStdout(t, func() {
		PrintFeedPost(post, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintFeedPostSingleJSON(t *testing.T) {
	post := FeedPostRow{ID: "p1", Author: "trader1", Message: "Test"}

	out := captureStdout(t, func() {
		PrintFeedPost(post, JSON)
	})

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}
