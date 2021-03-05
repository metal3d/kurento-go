package kurento

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/websocket"
)

var ErrConnectionClosing = errors.New("kurento: websocket connection is closing")

// Error that can be filled in response
type Error struct {
	Code    int64
	Message string
	Data    map[string]interface{}
}

func (e *Error) Error() string {
	return e.Message
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
	clientId      float64
	clientMap     threadsafeClientMap
	subscriberMap threadsafeSubscriberMap
	host          string
	ws            *websocket.Conn
	SessionId     string
	closeSig      chan bool
	closeWg       *sync.WaitGroup
	toCloseCh     chan string
}

type threadsafeClientMap struct {
	clients map[float64]chan Response
	lock    sync.RWMutex
}

type threadsafeSubscriberMap struct {
	subscribers map[string]map[string]chan Response
	lock        sync.RWMutex
}

func NewConnection(host string) (*Connection, error) {
	c := new(Connection)

	c.closeSig = make(chan bool)
	c.closeWg = &sync.WaitGroup{}
	c.toCloseCh = make(chan string, 1)
	go awaitClose(c.toCloseCh, c.closeSig, c.closeWg)

	c.clientMap = threadsafeClientMap{
		clients: make(map[float64]chan Response),
		lock:    sync.RWMutex{},
	}
	var err error

	conf, err := websocket.NewConfig(host+"/kurento", "http://127.0.0.1")
	if err != nil {
		return nil, fmt.Errorf("kurento: error creating new config: %v", err)
	}
	conf.Dialer = &net.Dialer{Timeout: 5 * time.Second}
	c.ws, err = websocket.DialConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("kurento: error dialing: %v", err)
	}
	c.host = host
	c.closeWg.Add(1)
	go c.handleResponse()
	return c, nil
}

func awaitClose(msgs chan string, closeSig chan bool, wg *sync.WaitGroup) {
	msg := <-msgs
	log.Println(fmt.Errorf("kurento: websocket closing with err%s", msg))
	close(closeSig)
	wg.Wait()
}

func (c *Connection) Create(m IMediaObject, options map[string]interface{}) error {
	elem := &MediaObject{}
	elem.setConnection(c)
	return elem.Create(m, options)
}

func (c *Connection) Close() error {
	select {
	case c.toCloseCh <- "Close called":
	case <-c.closeSig:
	}
	c.closeWg.Wait()
	return nil
}

type respErr struct {
	r   Response
	err error
}

func (c *Connection) handleResponse() {
	defer c.closeWg.Done()
	defer c.ws.Close()
	var err error
	var r Response
	var retVal respErr
	resps := make(chan respErr, 1)
	for {
		go func() {
			resp := Response{}
			var err error
			if debug {
				var test string
				err = websocket.Message.Receive(c.ws, &test)
				log.Println(test)
				json.Unmarshal([]byte(test), &resp)
			} else {
				err = websocket.JSON.Receive(c.ws, &resp)
			}
			resps <- respErr{
				r:   resp,
				err: err,
			}
		}()
		select {
		case retVal = <-resps:
		case <-c.closeSig:
			return
		}
		r = retVal.r
		err = retVal.err

		if err != nil {
			c.toCloseCh <- "Error reading from websocket: " + err.Error()
			return
		}

		if r.Result.SessionId != "" && c.SessionId != r.Result.SessionId {
			if debug {
				log.Println("SESSIONID RETURNED")
			}
			c.SessionId = r.Result.SessionId
		}
		// if webscocket client exists, send response to the chanel
		c.clientMap.lock.RLock()
		c.subscriberMap.lock.RLock()
		if c.clientMap.clients[r.Id] != nil {
			c.closeWg.Add(1)
			go func(r Response) {
				defer c.closeWg.Done()
				select {
				case c.clientMap.clients[r.Id] <- r:
				case <-c.closeSig:
				}
				// channel is read or we are closing, we can delete it
				c.clientMap.lock.Lock()
				close(c.clientMap.clients[r.Id])
				delete(c.clientMap.clients, r.Id)
				c.clientMap.lock.Unlock()
			}(r)
		} else if r.Method == "onEvent" && c.subscriberMap.subscribers[r.Params.Value.Data.Type][r.Params.Value.Data.Source] != nil {
			// Need to send it to the channel created on subscription
			c.closeWg.Add(1)
			go func(r Response) {
				defer c.closeWg.Done()
				c.subscriberMap.lock.RLock()
				select {
				case c.subscriberMap.subscribers[r.Params.Value.Data.Type][r.Params.Value.Data.Source] <- r:
				case <-c.closeSig:
				}
				c.subscriberMap.lock.RUnlock()
			}(r)
		} else if debug {
			if r.Method == "" {
				log.Println("Dropped message because there is no client ", r.Id)
			} else {
				log.Println("Dropped message because there is no subscription", r.Params.Value.Data.Type)
			}
			spew.Dump(r)
		}
		c.clientMap.lock.RUnlock()
		c.subscriberMap.lock.RUnlock()
	}
}

// Allow clients to subscribe to messages intended for them
func (c *Connection) Subscribe(eventType string, elementId string) <-chan Response {
	c.subscriberMap.lock.Lock()
	defer c.subscriberMap.lock.Unlock()
	if c.subscriberMap.subscribers == nil {
		c.subscriberMap.subscribers = make(map[string]map[string]chan Response)
	}
	if _, ok := c.subscriberMap.subscribers[eventType]; !ok {
		c.subscriberMap.subscribers[eventType] = make(map[string]chan Response)
	}
	c.subscriberMap.subscribers[eventType][elementId] = make(chan Response)
	return c.subscriberMap.subscribers[eventType][elementId]
}

// Allow clients to unsubscribe from messages intended for them
func (c *Connection) Unsubscribe(eventType string, elementId string) {
	c.subscriberMap.lock.Lock()
	defer c.subscriberMap.lock.Unlock()
	close(c.subscriberMap.subscribers[eventType][elementId])
	delete(c.subscriberMap.subscribers[eventType], elementId)
}

func (c *Connection) Request(req map[string]interface{}) (<-chan Response, error) {
	c.clientId++
	req["id"] = c.clientId
	if c.SessionId != "" {
		req["sessionId"] = c.SessionId
	}
	c.clientMap.lock.Lock()
	defer c.clientMap.lock.Unlock()
	c.clientMap.clients[c.clientId] = make(chan Response)
	if debug {
		j, _ := json.MarshalIndent(req, "", "    ")
		log.Println("json", string(j))
	}
	err := websocket.JSON.Send(c.ws, req)
	return c.clientMap.clients[c.clientId], err
}
