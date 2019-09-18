package kurento

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

// Base for all objects that can be created in the media server.
type MediaObject struct {
	connection *Connection

	// `MediaPipeline` to which this MediaObject belong, or the pipeline itself if
	// invoked over a `MediaPipeline`
	MediaPipeline IMediaPipeline

	// parent of this media object. The type of the parent depends on the type of the
	// element. The parent of a `MediaPad` is its `MediaElement`; the parent of a
	// `Hub` or a `MediaElement` is its `MediaPipeline`. A `MediaPipeline` has no
	// parent, i.e. the property is null
	Parent IMediaObject

	// unique identifier of the mediaobject.
	Id string

	// Childs of current object, all returned objects have parent set to current
	// object
	Childs []IMediaObject

	// Object name. This is just a comodity to simplify developers life debugging, it
	// is not used internally for indexing nor idenfiying the objects. By default is
	// the object type followed by the object id.
	Name string
}

// Return contructor params to be called by "Create".
func (elem *MediaObject) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

func (elem *MediaObject) Subscribe(eventType string, handler SubscriptionHandler) error {

	req := elem.getSubscribeRequest()

	// params := make(map[string]interface{})

	req["params"] = map[string]interface{}{
		"type":   eventType,
		"object": elem.Id,
	}
	// Call server and go run to the
	message := <-elem.connection.Request(req)
	if message.Error != nil {
		log.Println("Error trying to subscribe to " + eventType)
		return errors.New(fmt.Sprintf("[%d] %s %s", message.Error.Code, message.Error.Message, message.Error.Data))
	}

	c := elem.connection.Subscribe(eventType, elem.Id)
	go func() {
		for {
			msg := <-c
			handler.Handle(msg)
		}
	}()

	// Returns error or nil
	if message.Error != nil {
		return errors.New(fmt.Sprintf("[%d] %s %s", message.Error.Code, message.Error.Message, message.Error.Data))
	}
	return nil
}
func (elem *MediaObject) Unsubscribe(eventType string, subscriptionId string) error {
	elem.connection.Unsubscribe(eventType, elem.Id)
	return nil
}

func (elem *MediaObject) Release() error {

	req := elem.getReleaseRequest()

	// params := make(map[string]interface{})

	req["params"] = map[string]interface{}{
		"object": elem.Id,
	}
	// Call server and go run to the
	message := <-elem.connection.Request(req)
	if message.Error != nil {
		log.Println("Error trying to release " + elem.Id)
		return errors.New(fmt.Sprintf("[%d] %s %s", message.Error.Code, message.Error.Message, message.Error.Data))
	}
	if debug {
		log.Println("Release response: ", message)
	}

	// Returns error or nil
	if message.Error != nil {
		return errors.New(fmt.Sprintf("[%d] %s %s", message.Error.Code, message.Error.Message, message.Error.Data))
	}
	return nil
}

type IServerManager interface {
}

// This is a standalone object for managing the MediaServer
type ServerManager struct {
	MediaObject

	// Server information, version, modules, factories, etc
	Info *ServerInfo

	// All the pipelines available in the server
	Pipelines []IMediaPipeline

	// All active sessions in the server
	Sessions []string
}

// Return contructor params to be called by "Create".
func (elem *ServerManager) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type ISessionEndpoint interface {
}

// Session based endpoint. A session is considered to be started when the media
// exchange starts. On the other hand, sessions terminate when a timeout,
// defined by the developer, takes place after the connection is lost.
type SessionEndpoint struct {
	Endpoint
}

// Return contructor params to be called by "Create".
func (elem *SessionEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IHub interface {
}

// A Hub is a routing `MediaObject`. It connects several `endpoints <Endpoint>`
// together
type Hub struct {
	MediaObject
}

// Return contructor params to be called by "Create".
func (elem *Hub) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IFilter interface {
}

// Base interface for all filters. This is a certain type of `MediaElement`, that
// processes media injected through its sinks, and delivers the outcome through
// its sources.
type Filter struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *Filter) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IEndpoint interface {
}

// Base interface for all end points. An Endpoint is a `MediaElement`
// that allow `KMS` to interchange media contents with external systems,
// supporting different transport protocols and mechanisms, such as `RTP`,
// `WebRTC`, `HTTP`, "file:/" URLs... An "Endpoint" may
// contain both sources and sinks for different media types, to provide
// bidirectional communication.
type Endpoint struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *Endpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IHubPort interface {
}

// This `MediaElement` specifies a connection with a `Hub`
type HubPort struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *HubPort) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"hub": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

type IPassThrough interface {
}

// This `MediaElement` that just passes media through
type PassThrough struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *PassThrough) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

type IUriEndpoint interface {
	Pause() error
	Stop() error
}

// Interface for endpoints the require a URI to work. An example of this, would be
// a `PlayerEndpoint` whose URI property could be used to locate a file to stream
type UriEndpoint struct {
	Endpoint

	// The uri for this endpoint.
	Uri string
}

// Return contructor params to be called by "Create".
func (elem *UriEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Pauses the feed
func (elem *UriEndpoint) Pause() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "pause",
		"object":    elem.Id,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil {
		return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return nil

}

// Stops the feed
func (elem *UriEndpoint) Stop() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "stop",
		"object":    elem.Id,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil {
		return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return nil

}

type IMediaPipeline interface {
}

// A pipeline is a container for a collection of `MediaElements<MediaElement>` and
// `MediaMixers<MediaMixer>`. It offers the methods needed to control the
// creation and connection of elements inside a certain pipeline.
type MediaPipeline struct {
	MediaObject
}

// Return contructor params to be called by "Create".
func (elem *MediaPipeline) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options
}

// Set if latency stats are being collected for pipeline
func (elem *MediaObject) SetLatencyStats(value bool) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "latencyStats", value)

	req["params"] = map[string]interface{}{
		"operation":       "setLatencyStats",
		"object":          elem.Id,
		"operationParams": params,
	}
	log.Println(req)

	// Call server and wait response
	response := <-elem.connection.Request(req)

	if response.Error != nil {
		return errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return nil
}

type ISdpEndpoint interface {
	GenerateOffer() (string, error)
	ProcessOffer(offer string) (string, error)
	ProcessAnswer(answer string) (string, error)
	GetLocalSessionDescriptor() (string, error)
	GetRemoteSessionDescriptor() (string, error)
}

// Implements an SDP negotiation endpoint able to generate and process
// offers/responses and that configures resources according to
// negotiated Session Description
type SdpEndpoint struct {
	SessionEndpoint

	// Maximum video bandwidth for receiving.
	// Unit: kbps(kilobits per second).
	// 0: unlimited.
	// Default value: 500
	MaxVideoRecvBandwidth int
}

// Return contructor params to be called by "Create".
func (elem *SdpEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Request a SessionSpec offer.
// This can be used to initiate a connection.
// Returns:
// // The SDP offer.
func (elem *SdpEndpoint) GenerateOffer() (string, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "generateOffer",
		"object":    elem.Id,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)
	// fmt.Println(response.Result["value"])
	// // The SDP offer.
	if response.Error != nil {
		return "", errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return trimQuotes(string(response.Result.Value)), nil

}

// Request the NetworkConnection to process the given SessionSpec offer (from the
// remote User Agent)
// Returns:
// // The chosen configuration from the ones stated in the SDP offer
func (elem *SdpEndpoint) ProcessOffer(offer string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "offer", offer)

	req["params"] = map[string]interface{}{
		"operation":       "processOffer",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The chosen configuration from the ones stated in the SDP offer
	if response.Error != nil {
		return "", errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return trimQuotes(string(response.Result.Value)), nil

}

// Request the NetworkConnection to process the given SessionSpec answer (from the
// remote User Agent).
// Returns:
// // Updated SDP offer, based on the answer received.
func (elem *SdpEndpoint) ProcessAnswer(answer string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "answer", answer)

	req["params"] = map[string]interface{}{
		"operation":       "processAnswer",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // Updated SDP offer, based on the answer received.
	if response.Error != nil {
		return "", errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return trimQuotes(string(response.Result.Value)), nil

}

// This method gives access to the SessionSpec offered by this NetworkConnection.
// .. note:: This method returns the local MediaSpec, negotiated or not. If no
// offer has been generated yet, it returns null. It an offer has been
// generated it returns the offer and if an answer has been processed
// it returns the negotiated local SessionSpec.
// Returns:
// // The last agreed SessionSpec
func (elem *SdpEndpoint) GetLocalSessionDescriptor() (string, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "getLocalSessionDescriptor",
		"object":    elem.Id,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The last agreed SessionSpec
	if response.Error != nil {
		return "", errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return trimQuotes(string(response.Result.Value)), nil

}

// This method gives access to the remote session description.
// .. note:: This method returns the media previously agreed after a complete
// offer-answer exchange. If no media has been agreed yet, it returns null.
// Returns:
// // The last agreed User Agent session description
func (elem *SdpEndpoint) GetRemoteSessionDescriptor() (string, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "getRemoteSessionDescriptor",
		"object":    elem.Id,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The last agreed User Agent session description
	if response.Error != nil {
		return "", errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return trimQuotes(string(response.Result.Value)), nil

}

type IBaseRtpEndpoint interface {
}

// Base class to manage common RTP features.
type BaseRtpEndpoint struct {
	SdpEndpoint

	// Minimum video bandwidth for sending.
	// Unit: kbps(kilobits per second).
	// 0: unlimited.
	// Default value: 100
	MinVideoSendBandwidth int

	// Maximum video bandwidth for sending.
	// Unit: kbps(kilobits per second).
	// 0: unlimited.
	// Default value: 500
	MaxVideoSendBandwidth int
}

// Return contructor params to be called by "Create".
func (elem *BaseRtpEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IMediaElement interface {
	GetSourceConnections(mediaType MediaType, description string) (*[]ElementConnectionData, error)
	GetSinkConnections(mediaType MediaType, description string) (*[]ElementConnectionData, error)
	Connect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error
	Disconnect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error
	SetAudioFormat(caps AudioCaps) error
	SetVideoFormat(caps VideoCaps) error
	GetStats() (map[string]ElementStats, error)
}

// Basic building blocks of the media server, that can be interconnected through
// the API. A `MediaElement` is a module that encapsulates a specific media
// capability. They can be connected to create media pipelines where those
// capabilities are applied, in sequence, to the stream going through the
// pipeline.
// `MediaElement` objects are classified by its supported media type (audio,
// video, etc.)
type MediaElement struct {
	MediaObject
}

// Return contructor params to be called by "Create".
func (elem *MediaElement) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Get the connections information of the elements that are sending media to this
// element `MediaElement`
// Returns:
// // A list of the connections information that are sending media to this
// element.
// // The list will be empty if no sources are found.
func (elem *MediaElement) GetSourceConnections(mediaType MediaType, description string) (*[]ElementConnectionData, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "description", description)

	req["params"] = map[string]interface{}{
		"operation":       "getSourceConnections",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // A list of the connections information that are sending media to this
	// element.
	// // The list will be empty if no sources are found.

	ret := []ElementConnectionData{}
	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return &ret, nil

}

// Returns a list of the connections information of the elements that ere
// receiving media from this element.
// Returns:
// // A list of the connections information that arereceiving media from this
// // element. The list will be empty if no sinks are found.
func (elem *MediaElement) GetSinkConnections(mediaType MediaType, description string) (*[]ElementConnectionData, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "description", description)

	req["params"] = map[string]interface{}{
		"operation":       "getSinkConnections",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // A list of the connections information that arereceiving media from this
	// // element. The list will be empty if no sinks are found.

	ret := []ElementConnectionData{}
	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	return &ret, nil

}

// Connects two elements, with the given restrictions, current `MediaElement` will
// start emmit media to sink element. Connection could take place in the future,
// when both media element show capabilities for connecting with the given
// restrictions
func (elem *MediaElement) Connect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "sink", sink)
	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)
	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

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

// Disconnects two elements, with the given restrictions, current `MediaElement`
// stops sending media to sink element. If the previously requested connection
// didn't took place it is also removed
func (elem *MediaElement) Disconnect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "sink", sink)
	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)
	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

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

// Sets the type of data for the audio stream. MediaElements that do not support
// configuration of audio capabilities will raise an exception
func (elem *MediaElement) SetAudioFormat(caps AudioCaps) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "caps", caps)

	req["params"] = map[string]interface{}{
		"operation":       "setAudioFormat",
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

// Sets the type of data for the video stream. MediaElements that do not support
// configuration of video capabilities will raise an exception
func (elem *MediaElement) SetVideoFormat(caps VideoCaps) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "caps", caps)

	req["params"] = map[string]interface{}{
		"operation":       "setVideoFormat",
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

// Get the stats associated with this media element
func (elem *MediaElement) GetStats() (map[string]ElementStats, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	req["params"] = map[string]interface{}{
		"operation":       "getStats",
		"object":          elem.Id,
		"operationParams": params,
	}

	// Call server and wait response
	response := <-elem.connection.Request(req)
	// Check for error first
	if response.Error != nil {
		return nil, errors.New(fmt.Sprintf("[%d] %s %s", response.Error.Code, response.Error.Message, response.Error.Data))
	}
	// Otherwise should try to get stuff
	elementStats := make(map[string]ElementStats)
	if string(response.Result.Value) != "" {
		// Get the first element of the map and coerce it to an "ElementStats" object
		err := json.Unmarshal(response.Result.Value, &elementStats)
		if err != nil {
			log.Println(err)
			return nil, err
		}

	} else {
		return nil, errors.New("getStats Response Value was nil")
	}
	// Returns error or nil
	return elementStats, nil
}

func trimQuotes(v string) string {
	return strings.Trim(v, "\"")
}
