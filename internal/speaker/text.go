package speaker

import (
	"iter"
	"regexp"
	"strings"
)

type TextChunker interface {
	Chunk(text string) iter.Seq[string]
}

type LineBreakTextChunker struct{}

func (c *LineBreakTextChunker) Chunk(text string) iter.Seq[string] {
	lb := regexp.MustCompile("(\r\n|\n|\r)")
	lines := lb.Split(text, -1)
	return func(yield func(string) bool) {
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if !yield(trimmed) {
				return
			}
		}
	}
}

// Ensure TextChunker implementation
var _ TextChunker = new(LineBreakTextChunker)
