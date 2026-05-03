package speaker_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/heptaliane/katarive-voicevox-narrator-plugin/internal/speaker"
)

func TestLineBreakTextChunker(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		text     string
		expected []string
	}{
		"CR": {
			text:     "a\ra",
			expected: []string{"a", "a"},
		},
		"LF": {
			text:     "a\na",
			expected: []string{"a", "a"},
		},
		"CRLF": {
			text:     "a\r\na",
			expected: []string{"a", "a"},
		},
		"multiple lines": {
			text:     "a\r\n\n\ra",
			expected: []string{"a", "", "", "a"},
		},
		"space wrapped": {
			text:     " a a ",
			expected: []string{"a a"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			chunker := new(speaker.LineBreakTextChunker)
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
