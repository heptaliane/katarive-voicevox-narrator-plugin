package speaker_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/heptaliane/katarive-voicevox-narrator-plugin/internal/speaker"
)

func TestDelimiterTextChunker(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		text      string
		delimiter string
		expected  []string
	}{
		"CR": {
			text:      "a\r\na\na\ra",
			delimiter: "\r",
			expected:  []string{"a", "a\na", "a"},
		},
		"LF": {
			text:      "a\r\na\na\ra",
			delimiter: "\n",
			expected:  []string{"a", "a", "a\ra"},
		},
		"CRLF": {
			text:      "a\r\na\na\ra",
			delimiter: "\r\n",
			expected:  []string{"a", "a\na\ra"},
		},
		"space wrapped": {
			text:      " a a ",
			delimiter: "\n",
			expected:  []string{"a a"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			chunker := speaker.DelimiterTextChunker{
				Delimiter: tc.delimiter,
			}
			var chunks []string
			for chunk := range chunker.Chunk(tc.text) {
				chunks = append(chunks, chunk)
			}

			if diff := cmp.Diff(tc.expected, chunks); diff != "" {
				t.Errorf("Unmatched chunk result (got: +, want: -)\n %s", diff)
			}
		})
	}
}

func TestRecursiveTextChunker(t *testing.T) {
	t.Parallel()

	chunkers := []speaker.TextChunker{
		&speaker.DelimiterTextChunker{Delimiter: "\r\n"},
		&speaker.DelimiterTextChunker{Delimiter: "\n"},
		&speaker.DelimiterTextChunker{Delimiter: "\r"},
	}

	cases := map[string]struct {
		text      string
		expected  []string
		threshold int
	}{
		"CR": {
			text:      "a\ra",
			threshold: 3,
			expected:  []string{"a", "a"},
		},
		"LF": {
			text:      "a\na",
			threshold: 3,
			expected:  []string{"a", "a"},
		},
		"CRLF": {
			text:      "a\r\na",
			threshold: 3,
			expected:  []string{"a", "a"},
		},
		"multiple lines": {
			text:      "a\r\na\na\ra",
			threshold: 3,
			expected:  []string{"a", "a", "a", "a"},
		},
		"space wrapped": {
			text:      " a a ",
			threshold: 10,
			expected:  []string{"a a"},
		},
		"short lines": {
			text:      "a\na",
			threshold: 10,
			expected:  []string{"a\na"},
		},
		"non-ascii lines": {
			text:      "いろは\nにほへと",
			threshold: 9,
			expected:  []string{"いろは\nにほへと"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			chunker := speaker.NewRecursiveTextChunker(
				&speaker.RuneCountTextChunkTrigger{Threshold: tc.threshold},
				chunkers,
			)
			var chunks []string
			for chunk := range chunker.Chunk(tc.text) {
				chunks = append(chunks, chunk)
			}

			if diff := cmp.Diff(tc.expected, chunks); diff != "" {
				t.Errorf("Unmatched chunk result (got: +, want: -)\n %s", diff)
			}
		})
	}
}
