package errors

import (
	"fmt"
)

type UnsupportedSpeakerError struct {
	Name  string
	Style string
}

func (e *UnsupportedSpeakerError) Error() string {
	return fmt.Sprintf("No such speaker (name: '%s', style: '%s')", e.Name, e.Style)
}

var _ error = new(UnsupportedSpeakerError)

type VoiceVoxConnectionError struct {
	Body     []byte
	Endpoint string
}

func (e *VoiceVoxConnectionError) Error() string {
	return fmt.Sprintf(
		"Connection with VoiceVox failed (%s): %s",
		string(e.Endpoint),
		string(e.Body),
	)
}

var _ error = new(VoiceVoxConnectionError)
