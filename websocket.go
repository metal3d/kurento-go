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
}

type Connection struct {
	clientId  float64
	clients   map[float64]chan Response
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
		websocket.JSON.Receive(c.ws, &r)
		if r.Result["sessionId"] != "" {
			if debug {
				log.Println("SESSIONID RETURNED")
			}
			c.SessionId = r.Result["sessionId"]
		}
		// if webscocket client exists, send response to the chanel
		if c.clients[r.Id] != nil {
			c.clients[r.Id] <- r
			// chanel is read, we can delete it
			delete(c.clients, r.Id)
		} else if debug {
			log.Println("Dropped message because there is no client ", r.Id)
			log.Println(r)
		}

	}
}

func (c *Connection) Request(req map[string]interface{}) <-chan Response {
	c.clientId++
	req["id"] = c.clientId
	if c.SessionId != "" {
		req["sesionId"] = c.SessionId
	}
	c.clients[c.clientId] = make(chan Response)
	if debug {
		j, _ := json.MarshalIndent(req, "", "    ")
		log.Println("json", string(j))
	}
	websocket.JSON.Send(c.ws, req)
	return c.clients[c.clientId]
}
