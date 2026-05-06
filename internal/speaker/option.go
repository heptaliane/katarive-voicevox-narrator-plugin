package speaker

import (
	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"
)

const DEFAULT_SPEAKER_NAME string = "ずんだもん"
const DEFAULT_SPEAKER_STYLE string = "ノーマル"
const DEFAULT_SPEAKER_ID int = 3 // ずんだもん (ノーマル)

type narrateOptions struct {
	speakerName  string
	speakerStyle string
	speakerId    int
	encoding     pb.AudioEncoding
}

type NarrateOption func(opt *narrateOptions)

func WithSpeakerName(name string) NarrateOption {
	return func(opt *narrateOptions) {
		opt.speakerName = name
	}
}
func WithSpeakerStyle(style string) NarrateOption {
	return func(opt *narrateOptions) {
		opt.speakerStyle = style
	}
}
func WithSpeakerId(speakerId int) NarrateOption {
	return func(opt *narrateOptions) {
		opt.speakerId = speakerId
	}
}
func WithEncoding(encoding pb.AudioEncoding) NarrateOption {
	return func(opt *narrateOptions) {
		opt.encoding = encoding
	}
}

func newNarrateOptions() *narrateOptions {
	return &narrateOptions{
		speakerName:  DEFAULT_SPEAKER_NAME,
		speakerStyle: DEFAULT_SPEAKER_STYLE,
		speakerId:    DEFAULT_SPEAKER_ID,
	}
}
