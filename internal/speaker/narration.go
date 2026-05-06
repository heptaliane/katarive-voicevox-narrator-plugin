package speaker

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"
)

// =================================
// Interfaces for NarrationGenerator
// =================================

type NarrationGenerator interface {
	Do(ctx context.Context, path, text string, opts ...NarrateOption) error
	Speakers() []*pb.SpeakerInfo
}

// =================================
// NarrationGenerator implementation
// =================================

// -------------------------
// ChunkedNarrationGenerator
// -------------------------
type ChunkedNarrationGenerator struct {
	Chunker  TextChunker
	Handler  VoiceVoxHandler
	CacheDir string
}

func (g *ChunkedNarrationGenerator) Do(
	ctx context.Context,
	path, text string,
	opts ...NarrateOption,
) error {
	options := newNarrateOptions()
	for _, opt := range opts {
		opt(options)
	}

	basedir := filepath.Join(g.CacheDir, fmt.Sprintf("%03d", options.SpeakerId))
	os.MkdirAll(basedir, 0755)

	ext := audioExtension(options.Encoding)

	var ps []string
	chunks := g.Chunker.Chunk(text)
	for chunk := range chunks {
		p := filepath.Join(basedir, fmt.Sprintf("%s.%s", filename(chunk), ext))
		ps = append(ps, p)
		if _, err := os.Stat(p); err == nil {
			continue
		}

		audio, err := g.Handler.Narrate(ctx, chunk, options)
		if err != nil {
			return err
		}

		encoder := getEncoder(options.Encoding)
		audio, err = encoder(audio)
		if err != nil {
			return err
		}

		err = saveAudio(p, audio)
		if err != nil {
			return err
		}
	}
	return concatAudio(ps, path)
}
func (g *ChunkedNarrationGenerator) Speakers() []*pb.SpeakerInfo {
	var speakers []*pb.SpeakerInfo
	for _, s := range g.Handler.Speakers() {
		speakers = append(speakers, &pb.SpeakerInfo{
			Id:   int32(s.Id),
			Name: fmt.Sprintf("%s (%s)", s.Name, s.Style),
		})
	}
	return speakers
}

// Ensure NarrationGenerator implementation
var _ NarrationGenerator = new(ChunkedNarrationGenerator)

// -----------------
// Helper components
// -----------------
func filename(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}
func saveAudio(path string, audio []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(audio))
	return err
}
func concatAudio(files []string, dest string) error {
	lf, err := os.CreateTemp("", fmt.Sprintf("%s.txt", filename(dest)))
	if err != nil {
		return err
	}
	defer os.Remove(lf.Name())

	for _, path := range files {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		fmt.Fprintf(lf, "file '%s'\n", abs)
	}
	lf.Close()

	cmd := exec.Command(
		"ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", lf.Name(),
		"-c", "copy",
		dest,
	)
	if stdout, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg execution error: %v (%s)", err, stdout)
	}

	return nil
}
func audioExtension(encoding pb.AudioEncoding) string {
	switch encoding {
	case pb.AudioEncoding_AUDIO_ENCODING_MP3:
		return "mp3"
	case pb.AudioEncoding_AUDIO_ENCODING_M4A:
		return "m4a"
	}
	return "wav"
}
func getEncoder(encoding pb.AudioEncoding) AudioEncoder {
	switch encoding {
	case pb.AudioEncoding_AUDIO_ENCODING_MP3:
		return MP3Encoder
	case pb.AudioEncoding_AUDIO_ENCODING_M4A:
		return M4AEncoder
	}
	return NopEncoder
}
