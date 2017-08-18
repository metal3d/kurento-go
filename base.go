package kurento

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

var debug = false

// Debug activate debug information.
func SetDebug(state bool) {
	debug = state
}

type SubscriptionHandler interface {
	Handle(event Response)
}

// IMadiaElement implements some basic methods as getConstructorParams or Create().
type IMediaObject interface {

	// Return the constructor parameters
	getConstructorParams(IMediaObject, map[string]interface{}) map[string]interface{}

	// Each media object should be able to create another object
	// Those options are sent to getConstructorParams
	Create(IMediaObject, map[string]interface{})

	// Add a subscription to event of "type"
	Subscribe(eventType string, handler SubscriptionHandler) error

	// remove a subscription to event of
	Unsubscribe(eventType string, subscriptionId string) error

	// Release the underlying resources in kurento
	Release() error

	// Set ID of the element
	setId(string)

	//Implement Stringer
	String() string

	setParent(IMediaObject)
	addChild(IMediaObject)

	setConnection(*Connection)
}

// Create object "m" with given "options"
func (elem *MediaObject) Create(m IMediaObject, options map[string]interface{}) {
	req := elem.getCreateRequest()
	constparams := m.getConstructorParams(elem, options)
	// TODO params["sessionId"]
	req["params"] = map[string]interface{}{
		"type":              getMediaElementType(m),
		"constructorParams": constparams,
	}
	if debug {
		log.Printf("request to be sent: %+v\n", req)
	}

	m.setConnection(elem.connection)

	res := <-elem.connection.Request(req)

	if debug {
		log.Printf("Oncreate response: %+v\n", res)
		log.Println(len(res.Result.Value))
		if len(res.Result.Value) != 0 {
			log.Println(string(res.Result.Value))
		}
	}

	if len(res.Result.Value) != 0 {
		elem.addChild(m)
		//m.setParent(elem)
		m.setId(trimQuotes(string(res.Result.Value)))
	}
}

// Create an object in memory that represents a remote object without creating it
func HydrateMediaObject(id string, parent IMediaObject, c *Connection, elem IMediaObject) error {
	elem.setConnection(c)
	elem.setId(id)
	if parent != nil {
		parent.addChild(elem)
	}
	return nil
}

// Implement setConnection that allows element to handle connection
func (elem *MediaObject) setConnection(c *Connection) {
	elem.connection = c
}

// Set parent of current element
// BUG(recursion) a recursion happends while testing, I must find why
func (elem *MediaObject) setParent(m IMediaObject) {
	elem.Parent = m
}

// Append child to the element
func (elem *MediaObject) addChild(m IMediaObject) {
	elem.Childs = append(elem.Childs, m)
}

// setId set object id from a KMS response
func (m *MediaObject) setId(id string) {
	m.Id = id
}

// Build a prepared create request
func (m *MediaObject) getCreateRequest() map[string]interface{} {

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "create",
		"params":  make(map[string]interface{}),
	}
}

// Build a prepared invoke request
func (m *MediaObject) getInvokeRequest() map[string]interface{} {
	req := m.getCreateRequest()
	req["method"] = "invoke"

	return req
}

func (m *MediaObject) getSubscribeRequest() map[string]interface{} {
	req := m.getCreateRequest()
	req["method"] = "subscribe"

	return req
}

func (m *MediaObject) getReleaseRequest() map[string]interface{} {
	req := m.getCreateRequest()
	req["method"] = "release"

	return req
}

// String implements fmt.Stringer interface, return ID
func (m *MediaObject) String() string {
	return m.Id
}

// Return name of the object
func getMediaElementType(i interface{}) string {
	n := reflect.TypeOf(i).String()
	p := strings.Split(n, ".")
	return p[len(p)-1]
}

func mergeOptions(a, b map[string]interface{}) {
	for key, val := range b {
		a[key] = val
	}
}

func setIfNotEmpty(param map[string]interface{}, name string, t interface{}) {
	log.Println("in Set if not empty")
	log.Println(t)
	switch v := t.(type) {
	case string:
		if v != "" {
			log.Println("Set Map Type String")
			param[name] = v
		}
	case int, float64:
		if v != 0 {
			log.Println("Set Map Type num")
			param[name] = v
		}
	case bool:
		if v {
			log.Println("Set Map Type bool")
			param[name] = v
		}
	case IMediaObject, fmt.Stringer:
		if v != nil {
			val := fmt.Sprintf("%s", v)
			if val != "" {
				log.Println("Set Map Type Stringer")
				param[name] = val
			}
		}
	case IceCandidate:
		val := fmt.Sprintf("%s", v)
		if val != "" {
			log.Println("Set Map Type Candidate")
			param[name] = v
		}
	default:
		log.Println("Couldn't set map type")
		log.Println(v)
	}
}
