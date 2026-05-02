package speaker_test

import (
	"context"
	"testing"

	"github.com/heptaliane/katarive-voicevox-narrator-plugin/internal/errors"
	"github.com/heptaliane/katarive-voicevox-narrator-plugin/internal/speaker"
)

func TestHttpVoiceVoxHandlerNarrate(t *testing.T) {
	t.Parallel()

	server := "http://localhost:50021"

	ctx := context.Background()
	handler, err := speaker.NewHttpVoiceVoxHandler(ctx, server)
	if err != nil {
		t.Fatalf("Failed to initialize handler: %v", err)
	}

	cases := map[string]struct {
		options       []speaker.VoiceVoxOption
		expectedError error
	}{
		"default": {
			options: []speaker.VoiceVoxOption{},
		},
		"with style": {
			options: []speaker.VoiceVoxOption{
				speaker.WithSpeakerStyle("ささやき"),
			},
		},
		"with speaker": {
			options: []speaker.VoiceVoxOption{
				speaker.WithSpeakerName("四国めたん"),
			},
		},
		"invalid style": {
			options: []speaker.VoiceVoxOption{
				speaker.WithSpeakerStyle("invalid style"),
			},
			expectedError: &errors.UnsupportedSpeakerError{
				Style: "invalid style",
				Name:  speaker.DEFAULT_SPEAKER_NAME,
			},
		},
		"invalid speaker": {
			options: []speaker.VoiceVoxOption{
				speaker.WithSpeakerName("invalid speaker"),
			},
			expectedError: &errors.UnsupportedSpeakerError{
				Style: speaker.DEFAULT_SPEAKER_STYLE,
				Name:  "invalid speaker",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			_, err := handler.Narrate(ctx, "unittest", tc.options...)

			if tc.expectedError == nil {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
			} else {
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("Unmatched error: expected %v but got %v", tc.expectedError, err)
					return
				}
			}
		})
	}
}
