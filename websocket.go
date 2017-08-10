package kurento

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

// Error that can be filled in response
type Error struct {
	Code    int64
	Message string
	Data    string
}

// Implements error built-in interface
func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s %s", e.Code, e.Message, e.Data)
}

// Response represents server response
type Response struct {
	Jsonrpc string
	Id      float64
	Result  map[string]string // should change if result has no several form
	Error   *Error
	Method  string
	Params  params
}

type params struct {
	Value  map[string]interface{}
	Object string
	Type string
}

type value struct {
	Data data
}

type data struct {
	Candidate IceCandidate
	Source string
	Tags []string
	Timestamp string
	Type string
}

type Connection struct {
	clientId  float64
	clients   map[float64]chan Response
	subscribers map[string]map[string]chan Response
	host      string
	ws        *websocket.Conn
	SessionId string
}

var connections = make(map[string]*Connection)

func NewConnection(host string) *Connection {
	if connections[host] != nil {
		return connections[host]
	}

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

func (c *Connection) handleResponse() {
	for { // run forever
		r := Response{}
		var test string
		if debug {
			websocket.Message.Receive(c.ws, &test)
			log.Println(test)
			json.Unmarshal([]byte(test),&r)
		} else {
			websocket.JSON.Receive(c.ws, &r)
		}
		
		if r.Result["sessionId"] != "" {
			if debug {
				log.Println("SESSIONID RETURNED")
			}
			c.SessionId = r.Result["sessionId"]
		}

		var data map[string]interface{}

		if r.Params.Value["data"] != nil {
			log.Println(r.Params.Value["data"])
			data = r.Params.Value["data"].(map[string]interface{})
			log.Println(data)
		}

		// if webscocket client exists, send response to the channel
		if c.clients[r.Id] != nil {
			c.clients[r.Id] <- r
			// channel is read, we can delete it
			delete(c.clients, r.Id)
		} else if r.Method == "onEvent" && c.subscribers[data["type"].(string)][data["source"].(string)] != nil{
			// Need to send it to the channel created on subscription
			go func() {
				c.subscribers[data["type"].(string)][data["source"].(string)] <- r
			}()

		} else if debug {
			if r.Method == "" {
				log.Println("Dropped message because there is no client ", r.Id)
			} else {
				log.Println("Dropped message because there is no subscription", r.Params.Value["data"].(map[string]string)["type"])
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
	if _, ok := c.subscribers[eventType] ; !ok {
		c.subscribers[eventType] = make(map[string]chan Response)
	}
	c.subscribers[eventType][elementId] = make(chan Response)
	return c.subscribers[eventType][elementId]
}

// Allow clients to unsubscribe from messages intended for them
func (c *Connection) Unsubscribe(eventType string, elementId string) {
	delete(c.subscribers[eventType],elementId)
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
