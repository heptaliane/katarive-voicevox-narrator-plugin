package speaker

import (
	"iter"
	"regexp"
	"strings"
	"unicode/utf8"
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

type DelimiterTextChunker struct {
	Delimiter string
}

func (c *DelimiterTextChunker) Chunk(text string) iter.Seq[string] {
	lines := strings.Split(text, c.Delimiter)
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
var _ TextChunker = new(DelimiterTextChunker)

type RecursiveTextChunker struct {
	trigger TextChunkTrigger
	inner   TextChunker
	next    *RecursiveTextChunker
}

func (c *RecursiveTextChunker) Chunk(text string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for ip := range c.inner.Chunk(text) {
			if c.trigger.ShouldChunk(ip) && c.next != nil {
				for np := range c.next.Chunk(ip) {
					if !yield(np) {
						return
					}
				}
			} else {
				if !yield(ip) {
					return
				}
			}
		}
	}
}

// Ensure TextChunker implementation
var _ TextChunker = new(RecursiveTextChunker)

func NewRecursiveTextChunker(
	trigger TextChunkTrigger,
	chunkers []TextChunker,
) *RecursiveTextChunker {
	var rc *RecursiveTextChunker
	for i := len(chunkers) - 1; i >= 0; i -= 1 {
		rc = &RecursiveTextChunker{
			trigger: trigger,
			inner:   chunkers[i],
			next:    rc,
		}
	}
	return rc
}

// Helpers
type TextChunkTrigger interface {
	ShouldChunk(text string) bool
}

type RuneCountTextChunkTrigger struct {
	Threshold int
}

func (t *RuneCountTextChunkTrigger) ShouldChunk(text string) bool {
	return utf8.RuneCountInString(text) >= t.Threshold
}
