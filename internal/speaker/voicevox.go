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
	SpeakerId(name, style string) (int, error)
	Narrate(ctx context.Context, text string, speakerId int) ([]byte, error)
}

// ============================
// VoiceVoxHandler Implementation
// ============================

// ----------------------
// HttpVoiceVoxHandler
// ----------------------
type HttpVoiceVoxHandler struct {
	speakers []*voiceVoiceSpeaker
	client   voicevox.ClientWithResponsesInterface
}

func (h *HttpVoiceVoxHandler) Narrate(
	ctx context.Context,
	text string,
	speakerId int,
) ([]byte, error) {
	aqp := &voicevox.AudioQueryParams{
		Text:    text,
		Speaker: speakerId,
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

	sp := &voicevox.SynthesisParams{
		Speaker: speakerId,
	}
	res, err := h.client.SynthesisWithResponse(ctx, sp, *aq.JSON200)
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
func (h *HttpVoiceVoxHandler) SpeakerId(name, style string) (int, error) {
	for _, speaker := range h.speakers {
		if name == speaker.name && style == speaker.style {
			return speaker.id, nil
		}
	}
	return 0, &errors.UnsupportedSpeakerError{
		Name:  name,
		Style: style,
	}
}

// Ensure VoiceVoxHandler implementation
var _ VoiceVoxHandler = new(HttpVoiceVoxHandler)

// -----------------
// Helper components
// -----------------
type voiceVoiceSpeaker struct {
	name  string
	style string
	id    int
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

	var speakers []*voiceVoiceSpeaker
	for _, speaker := range *res.JSON200 {
		for _, style := range speaker.Styles {
			speakers = append(speakers, &voiceVoiceSpeaker{
				name:  speaker.Name,
				style: style.Name,
				id:    style.Id,
			})
		}
	}

	return &HttpVoiceVoxHandler{
		client:   client,
		speakers: speakers,
	}, nil
}
