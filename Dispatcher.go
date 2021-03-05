package kurento

import (
	"errors"
	"fmt"
)

type IDispatcher interface {
	Connect(source HubPort, sink HubPort) error
}

// A `Hub` that allows routing between arbitrary port pairs
type Dispatcher struct {
	Hub
}

// Return contructor params to be called by "Create".
func (elem *Dispatcher) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

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
func (elem *Dispatcher) Connect(source HubPort, sink HubPort) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "source", source)
	setIfNotEmpty(params, "sink", sink)

	req["params"] = map[string]interface{}{
		"operation":       "connect",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	responses, err := elem.connection.Request(req)
	if err != nil {
		return err
	}
	select {
	case response := <-responses:
		// Returns error or nil
		if response.Error != nil {
			return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
		}
	case <-elem.connection.closeSig:
		return ErrConnectionClosing
	}

	return nil

}
