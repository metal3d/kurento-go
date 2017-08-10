package kurento

import (
	"fmt"
	"errors"
	)

type IMixer interface {
	Connect(media MediaType, source HubPort, sink HubPort) error
	Disconnect(media MediaType, source HubPort, sink HubPort) error
}

// A `Hub` that allows routing of video between arbitrary port pairs and mixing of
// audio among several ports
type Mixer struct {
	Hub
}

// Return contructor params to be called by "Create".
func (elem *Mixer) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Connects each corresponding :rom:enum:`MediaType` of the given source port with
// the sink port.
func (elem *Mixer) Connect(media MediaType, source HubPort, sink HubPort) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "media", media)
	setIfNotEmpty(params, "source", source)
	setIfNotEmpty(params, "sink", sink)

	req["params"] = map[string]interface{}{
		"operation":       "connect",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil { 
		return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return nil

}

// Disonnects each corresponding :rom:enum:`MediaType` of the given source port
// from the sink port.
func (elem *Mixer) Disconnect(media MediaType, source HubPort, sink HubPort) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "media", media)
	setIfNotEmpty(params, "source", source)
	setIfNotEmpty(params, "sink", sink)

	req["params"] = map[string]interface{}{
		"operation":       "disconnect",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil { 
		return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return nil

}
