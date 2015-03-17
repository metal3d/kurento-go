package kurento

import "fmt"

type IRtpEndpoint interface {
}

// Endpoint that provides bidirectional content delivery capabilities with remote
// networked peers through RTP protocol. An `RtpEndpoint` contains paired sink
// and source `MediaPad` for audio and video.
type RtpEndpoint struct {
	SdpEndpoint
}

// Return contructor params to be called by "Create".
func (elem *RtpEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}
