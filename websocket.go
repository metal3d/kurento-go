package kurento

import (
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

// Error that can be filled in response
type Error struct {
	Code    int64
	Message string
	Data    string
}

// Response represents server response
type Response struct {
	Jsonrpc string
	Id      float64
	Result  result // should change if result has no several form
	Error   *Error
	Method  string
	Params  Params
}
type result struct {
	Value     json.RawMessage
	SessionId string
	Object    string
}
type Params struct {
	Value  Value
	Object string
	Type   string
}

type Value struct {
	Data Data
}

type Data struct {
	Candidate IceCandidate
	Source    string
	Tags      []string
	Timestamp string
	Type      string
	State     string
	StreamId  int
}

type Connection struct {
	clientId    float64
	clients     map[float64]chan Response
	subscribers map[string]map[string]chan Response
	host        string
	ws          *websocket.Conn
	SessionId   string
}

var connections = make(map[string]*Connection)

func NewConnection(host string) *Connection {
	// if connections[host] != nil {
	// 	return connections[host]
	// }

	c := new(Connection)
	connections[host] = c

	c.clients = make(map[float64]chan Response)
	var err error
	c.ws, err = websocket.Dial(host+"/kurento", "", "http://127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	c.host = host
	go c.handleResponse()
	return c
}

func (c *Connection) Create(m IMediaObject, options map[string]interface{}) {
	elem := &MediaObject{}
	elem.setConnection(c)
	elem.Create(m, options)
}

func (c *Connection) Close() error {
	return c.ws.Close()
}

func (c *Connection) handleResponse() {
	var err error
	var test string
	var r Response
	for { // run forever
		r = Response{}
		if debug {
			err = websocket.Message.Receive(c.ws, &test)
			log.Println(test)
			json.Unmarshal([]byte(test), &r)
		} else {
			err = websocket.JSON.Receive(c.ws, &r)
		}
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
		}

		if r.Result.SessionId != "" && c.SessionId != r.Result.SessionId {
			if debug {
				log.Println("SESSIONID RETURNED")
			}
			c.SessionId = r.Result.SessionId
		}
		// if webscocket client exists, send response to the chanel
		if c.clients[r.Id] != nil {
			c.clients[r.Id] <- r
			// channel is read, we can delete it
			delete(c.clients, r.Id)
		} else if r.Method == "onEvent" && c.subscribers[r.Params.Value.Data.Type][r.Params.Value.Data.Source] != nil {
			// Need to send it to the channel created on subscription
			go func() {
				c.subscribers[r.Params.Value.Data.Type][r.Params.Value.Data.Source] <- r
			}()

		} else if debug {
			if r.Method == "" {
				log.Println("Dropped message because there is no client ", r.Id)
			} else {
				log.Println("Dropped message because there is no subscription", r.Params.Value.Data.Type)
			}
			log.Println(r)
		}

	}
}

// Allow clients to subscribe to messages intended for them
func (c *Connection) Subscribe(eventType string, elementId string) <-chan Response {
	if c.subscribers == nil {
		c.subscribers = make(map[string]map[string]chan Response)
	}
	if _, ok := c.subscribers[eventType]; !ok {
		c.subscribers[eventType] = make(map[string]chan Response)
	}
	c.subscribers[eventType][elementId] = make(chan Response)
	return c.subscribers[eventType][elementId]
}

// Allow clients to unsubscribe from messages intended for them
func (c *Connection) Unsubscribe(eventType string, elementId string) {
	delete(c.subscribers[eventType], elementId)
}

func (c *Connection) Request(req map[string]interface{}) <-chan Response {
	c.clientId++
	req["id"] = c.clientId
	if c.SessionId != "" {
		req["sessionId"] = c.SessionId
	}
	c.clients[c.clientId] = make(chan Response)
	if debug {
		j, _ := json.MarshalIndent(req, "", "    ")
		log.Println("json", string(j))
	}
	websocket.JSON.Send(c.ws, req)
	return c.clients[c.clientId]
}
