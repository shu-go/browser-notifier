package main

import (
	"log"
)

var _ = log.Print

type Server struct {
	clients []*Client

	sendChan         chan interface{}
	closeChan        chan struct{}
	clientAppendChan chan *Client
	clientRemoveChan chan *Client
}

func NewServer() *Server {
	return &Server{
		sendChan:         make(chan interface{}),
		closeChan:        make(chan struct{}),
		clientAppendChan: make(chan *Client),
		clientRemoveChan: make(chan *Client),
	}
}

func (s *Server) Send(v interface{}) {
	s.sendChan <- v
}

func (s *Server) Close() {
	s.closeChan <- struct{}{}
}

func (s *Server) AppendClient(c *Client) {
	s.clientAppendChan <- c
}

func (s *Server) RemoveClient(c *Client) {
	s.clientRemoveChan <- c
}

func (s *Server) Run() {
	for {
		select {
		case <-s.closeChan:
			s.doClose()

			return

		case v := <-s.sendChan:
			s.doSend(v)

		case c := <-s.clientAppendChan:
			s.doAppendClient(c)

		case c := <-s.clientRemoveChan:
			s.doRemoveClient(c)

		}
		//log.Println("Server: Loop")
	}
}

func (s *Server) doClose() {
	for _, c := range s.clients {
		// TODO: need timeout?
		c.Close()
	}
	s.clients = nil
}

func (s *Server) doSend(v interface{}) {
	ng := []*Client{}

	for _, c := range s.clients {
		// TODO: need timeout?
		err := c.SendSync(v)

		if err != nil {
			ng = append(ng, c)
		}
	}

	// disconnect & remove ng clients
	for _, ngc := range ng {
		s.doRemoveClient(ngc)
	}
}

func (s *Server) doAppendClient(c *Client) {
	s.clients = append(s.clients, c)
}

func (s *Server) doRemoveClient(c *Client) {
	idx := s.index(c)
	if idx == -1 {
		return
	}
	s.clients = append(s.clients[:idx], s.clients[idx+1:]...)
}

func (s *Server) index(c *Client) int {
	for i, ic := range s.clients {
		if ic == c || ic.id == c.id {
			return i
		}
	}
	return -1
}
