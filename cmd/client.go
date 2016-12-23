package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"sync"
	"time"
)

var _ = log.Print

type Client struct {
	id string
	ws *websocket.Conn

	sendChan  chan interface{}
	closeChan chan struct{}
}

var (
	mu     sync.Mutex
	lastID int64
)

func init() {
	lastID = time.Now().UnixNano()
}

//TODO: Receive() that is called in Run()
//type Receiver func(v interface{}, c *Client)

func NewClient(ws *websocket.Conn) *Client {
	mu.Lock()
	now := time.Now().UnixNano()
	if now == lastID {
		now += 1
	}
	lastID = now
	mu.Unlock()

	return &Client{id: fmt.Sprintf("%d", now), ws: ws}
}

func NewClientWithID(id string, ws *websocket.Conn) *Client {
	return &Client{id: id, ws: ws}
}

func (c *Client) String() string {
	return fmt.Sprintf("[%d] %s", c.id, c.ws.RemoteAddr())
}

// discard err!! If neeeded, use SendSync() instead.
func (c *Client) Send(v interface{}) error {
	c.sendChan <- v
	return nil
}

func (c *Client) SendSync(v interface{}) error {
	return c.doSend(v)
}

func (c *Client) Close() {
	c.closeChan <- struct{}{}
}

func (c *Client) Run() {
	for {
		select {
		case <-c.closeChan:
			c.doClose()

			return

		case v := <-c.sendChan:
			c.doSend(v)
		}
	}
	//log.Println("Client: Loop")
}

func (c *Client) doClose() {
	if c.ws != nil {
		c.ws.Close()
	}
	c.ws = nil
}

func (c *Client) doSend(v interface{}) error {
	return websocket.JSON.Send(c.ws, v)
}
