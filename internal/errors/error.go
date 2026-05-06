package errors

import (
	"fmt"
)

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
