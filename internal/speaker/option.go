package speaker

type NarrateOptions struct {
	speakerName  string
	speakerStyle string
}

type NarrateOption func(opt *NarrateOptions)

func WithSpeakerName(name string) NarrateOption {
	return func(opt *NarrateOptions) {
		opt.speakerName = name
	}
}
func WithSpeakerStyle(style string) NarrateOption {
	return func(opt *NarrateOptions) {
		opt.speakerStyle = style
	}
}

func newNarrateOptions() *NarrateOptions {
	return &NarrateOptions{
		speakerName:  DEFAULT_SPEAKER_NAME,
		speakerStyle: DEFAULT_SPEAKER_STYLE,
	}
}
