package speaker_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/heptaliane/katarive-voicevox-narrator-plugin/internal/speaker"
)

func TestChunkedNarrationGenerator(t *testing.T) {
	t.Parallel()

	server := "http://localhost:50021"
	basedir := t.TempDir()
	output := filepath.Join(basedir, "output.wav")

	ctx := context.Background()
	handler, err := speaker.NewHttpVoiceVoxHandler(ctx, server)
	if err != nil {
		t.Fatalf("Failed to initialize handler: %v", err)
	}
	generator := &speaker.ChunkedNarrationGenerator{
		Chunker:  new(speaker.LineBreakTextChunker),
		Handler:  handler,
		CacheDir: basedir,
	}

	err = generator.Do(
		ctx,
		output,
		`
		Hello
		World
		`,
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if _, err := os.Stat(output); err != nil {
		t.Errorf("Failed to create output file")
	}
}
