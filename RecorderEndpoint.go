package kurento

import "fmt"

type IRecorderEndpoint interface {
	Record() error
}

// Provides function to store contents in reliable mode (doesn't discard data). It
// contains `MediaSink` pads for audio and video.
type RecorderEndpoint struct {
	UriEndpoint
}

// Return contructor params to be called by "Create".
func (elem *RecorderEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":     fmt.Sprintf("%s", from),
		"uri":               "",
		"mediaProfile":      fmt.Sprintf("%s", from),
		"stopOnEndOfStream": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Starts storing media received through the `MediaSink` pad
func (elem *RecorderEndpoint) Record() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "record",
		"object":    elem.Id,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	return response.Error

}
