package output

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type FeedPostRow struct {
	ID       string
	Author   string
	Message  string
	Created  string
}

func PrintFeedPosts(posts []FeedPostRow, format Format) {
	if format == JSON {
		PrintJSON(posts)
		return
	}

	if len(posts) == 0 {
		Infof("No feed posts found.")
		return
	}

	t := NewTable("Author", "Message", "Posted")
	for _, p := range posts {
		msg := p.Message
		if len(msg) > 80 {
			msg = msg[:77] + "..."
		}
		msg = strings.ReplaceAll(msg, "\n", " ")
		t.AppendRow(table.Row{
			Cyan(p.Author),
			msg,
			p.Created,
		})
	}
	RenderTable(t)
}

func PrintPostCreated(format Format) {
	if format == JSON {
		PrintJSON(map[string]any{"status": "created"})
		return
	}

	Successf("Post created successfully.")
}

func PrintFeedPost(post FeedPostRow, format Format) {
	if format == JSON {
		PrintJSON(post)
		return
	}

	fmt.Printf("\n  %s  %s\n", Bold(post.Author), post.Created)
	fmt.Printf("  %s\n\n", post.Message)
}
