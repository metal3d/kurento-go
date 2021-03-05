package kurento

import (
	"errors"
	"fmt"
)

type IHttpGetEndpoint interface {
}

// An "HttpGetEndpoint" contains SOURCE pads for AUDIO and VIDEO, delivering media
// using HTML5 pseudo-streaming mechanism.
// This type of endpoint provide unidirectional communications. Its `MediaSink`
// is associated with the HTTP GET method
type HttpGetEndpoint struct {
	HttpEndpoint
}

// Return contructor params to be called by "Create".
func (elem *HttpGetEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":        fmt.Sprintf("%s", from),
		"terminateOnEOS":       fmt.Sprintf("%s", from),
		"mediaProfile":         fmt.Sprintf("%s", from),
		"disconnectionTimeout": 2,
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

type IHttpPostEndpoint interface {
}

// An `HttpPostEndpoint` contains SINK pads for AUDIO and VIDEO, which provide
// access to an HTTP file upload function
// This type of endpoint provide unidirectional communications. Its
// `MediaSources <MediaSource>` are accessed through the `HTTP` POST method.
type HttpPostEndpoint struct {
	HttpEndpoint
}

// Return contructor params to be called by "Create".
func (elem *HttpPostEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":        fmt.Sprintf("%s", from),
		"disconnectionTimeout": 2,
		"useEncodedMedia":      fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

type IHttpEndpoint interface {
	GetUrl() (string, error)
}

// Endpoint that enables Kurento to work as an HTTP server, allowing peer HTTP
// clients to access media.
type HttpEndpoint struct {
	SessionEndpoint
}

// Return contructor params to be called by "Create".
func (elem *HttpEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Obtains the URL associated to this endpoint
// Returns:
// // The url as a String
func (elem *HttpEndpoint) GetUrl() (string, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "getUrl",
		"object":    elem.Id,
	}

	// Call server and wait response
	responses, err := elem.connection.Request(req)
	if err != nil {
		return "", err
	}
	var response Response
	select {
	case response = <-responses:
	case <-elem.connection.closeSig:
		return "", ErrConnectionClosing
	}
	// // The url as a String
	if response.Error != nil {
		return "", errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return trimQuotes(string(response.Result.Value)), nil

}
