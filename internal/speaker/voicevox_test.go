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
		options       []speaker.NarrateOption
		expectedError error
	}{
		"default": {
			options: []speaker.NarrateOption{},
		},
		"with style": {
			options: []speaker.NarrateOption{
				speaker.WithSpeakerStyle("ささやき"),
			},
		},
		"with speaker": {
			options: []speaker.NarrateOption{
				speaker.WithSpeakerName("四国めたん"),
			},
		},
		"invalid style": {
			options: []speaker.NarrateOption{
				speaker.WithSpeakerStyle("invalid style"),
			},
			expectedError: &errors.UnsupportedSpeakerError{
				Style: "invalid style",
				Name:  speaker.DEFAULT_SPEAKER_NAME,
			},
		},
		"invalid speaker": {
			options: []speaker.NarrateOption{
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

			err := func() error {
				ctx := context.Background()
				id, err := handler.SpeakerId(tc.options...)
				if err != nil {
					return err
				}

				_, err = handler.Narrate(ctx, "unittest", id)
				return err
			}()

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
