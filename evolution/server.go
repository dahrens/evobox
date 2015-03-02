package evolution

import (
	"log"
)

// Chat server.
type Server struct {
	pattern  string
	messages []*Message
	clients  map[int]*Client
	addCh    chan *Client
	delCh    chan *Client
	doneCh   chan bool
	ErrCh    chan error
	Next     *AutoInc
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := make([]*Message, 100)
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	doneCh := make(chan bool)
	ErrCh := make(chan error)
	Next := NewAutoInc(0, 1)

	return &Server{
		pattern,
		messages,
		clients,
		addCh,
		delCh,
		doneCh,
		ErrCh,
		Next,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.ErrCh <- err
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	for {
		select {
		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.Id] = c
			log.Println("Now", len(s.clients), "clients connected.")

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.Id)

		// broadcast message for all clients
		// case msg := <-s.sendAllCh:
		// 	log.Println("Send all:", msg)
		// 	s.messages = append(s.messages, msg)
		// 	s.sendAll(msg)

		case err := <-s.ErrCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
