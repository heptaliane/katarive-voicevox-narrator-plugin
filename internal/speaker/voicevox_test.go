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
		options       speaker.NarrateOptions
		expectedError error
	}{
		"default": {
			options: speaker.NarrateOptions{
				SpeakerId: 3,
			},
		},
		"invalid id": {
			options: speaker.NarrateOptions{
				SpeakerId: -1,
			},
			expectedError: &errors.VoiceVoxConnectionError{
				Body:     []byte(`{"detail":"Internal Server Error"}`),
				Endpoint: "/audio_query",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := func() error {
				ctx := context.Background()
				if err != nil {
					return err
				}

				_, err = handler.Narrate(ctx, "unittest", &tc.options)
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
