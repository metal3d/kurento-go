package kurento

import (
	"errors"
	"fmt"
)

type IAlphaBlending interface {
	SetMaster(source HubPort, zOrder int) error
	SetPortProperties(relativeX float64, relativeY float64, zOrder int, relativeWidth float64, relativeHeight float64, port HubPort) error
}

// A `Hub` that mixes the :rom:attr:`MediaType.AUDIO` stream of its connected
// sources and constructs one output with :rom:attr:`MediaType.VIDEO`
// streams of its connected sources into its sink
type AlphaBlending struct {
	Hub
}

// Return contructor params to be called by "Create".
func (elem *AlphaBlending) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Sets the source port that will be the master entry to the mixer
func (elem *AlphaBlending) SetMaster(source HubPort, zOrder int) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "source", source)
	setIfNotEmpty(params, "zOrder", zOrder)

	req["params"] = map[string]interface{}{
		"operation":       "setMaster",
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
		if response.Error != nil {
			return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
		}
	case <-elem.connection.closeSig:
		return ErrConnectionClosing
	}

	return nil

}

// Configure the blending mode of one port.
func (elem *AlphaBlending) SetPortProperties(relativeX float64, relativeY float64, zOrder int, relativeWidth float64, relativeHeight float64, port HubPort) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "relativeX", relativeX)
	setIfNotEmpty(params, "relativeY", relativeY)
	setIfNotEmpty(params, "zOrder", zOrder)
	setIfNotEmpty(params, "relativeWidth", relativeWidth)
	setIfNotEmpty(params, "relativeHeight", relativeHeight)
	setIfNotEmpty(params, "port", port)

	req["params"] = map[string]interface{}{
		"operation":       "setPortProperties",
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
		if response.Error != nil {
			return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
		}
	case <-elem.connection.closeSig:
		return ErrConnectionClosing
	}

	return nil

}
