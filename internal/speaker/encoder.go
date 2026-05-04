package speaker

import (
	"bytes"
	"fmt"
	"os/exec"
)

type AudioEncoder = func(audio []byte) ([]byte, error)

func NopEncoder(audio []byte) ([]byte, error) {
	return audio, nil
}
func MP3Encoder(audio []byte) ([]byte, error) {
	return encodeWithFfmpeg(
		audio,
		"-f", "mp3", "-codec:a", "libmp3lame", "-qscale:a", "2",
	)
}
func M4AEncoder(audio []byte) ([]byte, error) {
	return encodeWithFfmpeg(
		audio,
		"-f", "ipod", "-codec:a", "aac", "-vbr", "3",
	)
}

// helpers
func encodeWithFfmpeg(audio []byte, ffmpegArgs ...string) ([]byte, error) {
	args := []string{"-i", "pipe:0"}
	args = append(args, ffmpegArgs...)
	args = append(args, "pipe:1")

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdin = bytes.NewReader(audio)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg conversion error: %v (%s)", err, stderr.String())
	}
	return stdout.Bytes(), nil
}
