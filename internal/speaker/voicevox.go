package speaker

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-hclog"

	"github.com/heptaliane/katarive-voicevox-narrator-plugin/gen/voicevox"
	"github.com/heptaliane/katarive-voicevox-narrator-plugin/internal/errors"
)

const DEFAULT_SPEAKER_NAME string = "ずんだもん"
const DEFAULT_SPEAKER_STYLE string = "ノーマル"

// ==============================
// Interfaces for VoiceVox handlers
// ==============================
type VoiceVoxHandler interface {
	Narrate(ctx context.Context, text string, opts ...VoiceVoxOption) ([]byte, error)
}

// -----------------
// Helper components
// -----------------

type voiceVoxOption struct {
	speakerName  string
	speakerStyle string
}

type VoiceVoxOption func(opt *voiceVoxOption)

func WithSpeakerName(name string) VoiceVoxOption {
	return func(opt *voiceVoxOption) {
		opt.speakerName = name
	}
}
func WithSpeakerStyle(style string) VoiceVoxOption {
	return func(opt *voiceVoxOption) {
		opt.speakerStyle = style
	}
}

func newVoiceVoxOption() *voiceVoxOption {
	return &voiceVoxOption{
		speakerName:  DEFAULT_SPEAKER_NAME,
		speakerStyle: DEFAULT_SPEAKER_STYLE,
	}
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

	logger hclog.Logger
}

func (h *HttpVoiceVoxHandler) Narrate(
	ctx context.Context,
	text string,
	opts ...VoiceVoxOption,
) ([]byte, error) {
	options := newVoiceVoxOption()
	for _, opt := range opts {
		opt(options)
	}

	speaker, err := h.getVoiceVoxId(options)
	if err != nil {
		return nil, err
	}

	aqp := &voicevox.AudioQueryParams{
		Text:    text,
		Speaker: speaker,
	}
	aq, err := h.client.AudioQueryWithResponse(ctx, aqp)
	if err != nil {
		return nil, err
	}
	if aq.StatusCode() != http.StatusOK {
		return nil, &errors.VoiceVoxConnectionError{Body: aq.Body}
	}

	sp := &voicevox.SynthesisParams{
		Speaker: speaker,
	}
	res, err := h.client.SynthesisWithResponse(ctx, sp, *aq.JSON200)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, &errors.VoiceVoxConnectionError{Body: res.Body}
	}

	return res.Body, nil
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

func (h *HttpVoiceVoxHandler) getVoiceVoxId(opts *voiceVoxOption) (int, error) {
	for _, speaker := range h.speakers {
		if opts.speakerName == speaker.name && opts.speakerStyle == speaker.style {
			return speaker.id, nil
		}
	}
	return 0, &errors.UnsupportedSpeakerError{
		Name:  opts.speakerName,
		Style: opts.speakerStyle,
	}
}

func NewHttpVoiceVoxHandler(
	ctx context.Context,
	server string,
	logger hclog.Logger,
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
		return nil, &errors.VoiceVoxConnectionError{Body: res.Body}
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
		logger:   logger,
	}, nil
}
