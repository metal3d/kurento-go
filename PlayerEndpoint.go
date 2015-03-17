package kurento

import "fmt"

type IPlayerEndpoint interface {
	Play() error
}

// Retrieves content from seekable sources in reliable
// mode (does not discard media information) and inject
// them into `KMS`. It
// contains one `MediaSource` for each media type detected.
type PlayerEndpoint struct {
	UriEndpoint
}

// Return contructor params to be called by "Create".
func (elem *PlayerEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":   fmt.Sprintf("%s", from),
		"uri":             "",
		"useEncodedMedia": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Starts to send data to the endpoint `MediaSource`
func (elem *PlayerEndpoint) Play() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "play",
		"object":    elem.Id,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	return response.Error

}
