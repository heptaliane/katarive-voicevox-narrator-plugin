package speaker

import (
	"context"
	"net/http"

	"github.com/heptaliane/katarive-voicevox-narrator-plugin/gen/voicevox"
	"github.com/heptaliane/katarive-voicevox-narrator-plugin/internal/errors"
)

// ==============================
// Interfaces for VoiceVox handlers
// ==============================
type VoiceVoxHandler interface {
	Narrate(ctx context.Context, text string, options *NarrateOptions) ([]byte, error)
	Speakers() []*VoiceVoxSpeaker
}

// ============================
// VoiceVoxHandler Implementation
// ============================

// ----------------------
// HttpVoiceVoxHandler
// ----------------------
type HttpVoiceVoxHandler struct {
	speakers []*VoiceVoxSpeaker
	client   voicevox.ClientWithResponsesInterface
}

func (h *HttpVoiceVoxHandler) Narrate(
	ctx context.Context,
	text string,
	options *NarrateOptions,
) ([]byte, error) {
	aqp := &voicevox.AudioQueryParams{
		Text:    text,
		Speaker: options.SpeakerId,
	}
	aq, err := h.client.AudioQueryWithResponse(ctx, aqp)
	if err != nil {
		return nil, err
	}
	if aq.StatusCode() != http.StatusOK {
		return nil, &errors.VoiceVoxConnectionError{
			Body:     aq.Body,
			Endpoint: "/audio_query",
		}
	}

	req := aq.JSON200
	req.PostPhonemeLength = options.PhonemeLength
	sp := &voicevox.SynthesisParams{
		Speaker: options.SpeakerId,
	}
	res, err := h.client.SynthesisWithResponse(ctx, sp, *req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, &errors.VoiceVoxConnectionError{
			Body:     res.Body,
			Endpoint: "/synthesis",
		}
	}

	return res.Body, nil
}
func (h *HttpVoiceVoxHandler) Speakers() []*VoiceVoxSpeaker {
	return h.speakers
}

// Ensure VoiceVoxHandler implementation
var _ VoiceVoxHandler = new(HttpVoiceVoxHandler)

// -----------------
// Helper components
// -----------------
type VoiceVoxSpeaker struct {
	Id    int
	Name  string
	Style string
}

func NewHttpVoiceVoxHandler(
	ctx context.Context,
	server string,
) (*HttpVoiceVoxHandler, error) {
	client, err := voicevox.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	sp := &voicevox.SpeakersParams{}
	res, err := client.SpeakersWithResponse(ctx, sp)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, &errors.VoiceVoxConnectionError{
			Body:     res.Body,
			Endpoint: "/speakers",
		}
	}

	var speakers []*VoiceVoxSpeaker
	for _, speaker := range *res.JSON200 {
		for _, style := range speaker.Styles {
			speakers = append(speakers, &VoiceVoxSpeaker{
				Name:  speaker.Name,
				Style: style.Name,
				Id:    style.Id,
			})
		}
	}

	return &HttpVoiceVoxHandler{
		client:   client,
		speakers: speakers,
	}, nil
}
