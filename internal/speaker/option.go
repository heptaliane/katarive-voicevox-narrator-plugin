package speaker

import (
	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"
)

const DEFAULT_SPEAKER_ID int = 3 // ずんだもん (ノーマル)
const DEFAULT_PHONEME_LENGTH float32 = 0.5

type NarrateOptions struct {
	SpeakerId     int
	PhonemeLength float32
	Encoding      pb.AudioEncoding
}

type NarrateOption func(opt *NarrateOptions)

func WithSpeakerId(speakerId int) NarrateOption {
	return func(opt *NarrateOptions) {
		opt.SpeakerId = speakerId
	}
}
func WithEncoding(encoding pb.AudioEncoding) NarrateOption {
	return func(opt *NarrateOptions) {
		opt.Encoding = encoding
	}
}
func WithPhonemeLength(length float32) NarrateOption {
	return func(opt *NarrateOptions) {
		opt.PhonemeLength = length
	}
}

func newNarrateOptions() *NarrateOptions {
	return &NarrateOptions{
		SpeakerId:     DEFAULT_SPEAKER_ID,
		PhonemeLength: DEFAULT_PHONEME_LENGTH,
	}
}
