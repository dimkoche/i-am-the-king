package main

import "net"
import "fmt"
import "log"
import "bufio"
//import "strings"
import "time"

type Server struct {
	priority string
	state    string
	king     string
	joins    chan net.Conn
	pong     bool
}

type Message struct {
	msg      string
	priority string
}

func (s *Server) Start() {
	ln, _ := net.Listen("tcp", PRIORITY_MAP[s.priority])

	go func() {
		for {
			select {
			case conn := <-s.joins:
				s.handleRequest(conn)
			}
		}
	}()

	go func() {
		for {
			conn, _ := ln.Accept()
			s.joins <- conn
		}
	}()

	go func() {
		c := time.Tick(2*T)
		for _ = range c {
			if s.state == "worker" {
				go s.pingKing()
			}
		}
	}()

	s.startElections()
}

func (s *Server) pingKing() {
	fmt.Println("pingKing")
	if s.king == "" || s.state == "king" {
		return
	}

	if s.pong {
		s.sendMessage("PING", PRIORITY_MAP[s.king])
		s.pong = false
	} else {
		go s.startElections()
	}

	// finish := time.After(time.Duration(4*T))
	// for {
	// 	select {
	// 	case <-finish:
	// 		fmt.Println("Waiting for king timeout")
	// 		if s.state == "waiting-for-king" {
	// 			//go s.startElections()
	// 			fmt.Println("Election should be started")
	// 		}
	// 	}
	// }
}

func (s *Server) startElections() {
	log.Printf("Start elections")
	s.setState("election")
	for priority, server := range PRIORITY_MAP {
		if s.priority < priority {
			go s.sendMessage("ALIVE?", server)
		}
	}

	finish := time.After(time.Duration(T))
	for {
		select {
		case <-finish:
			fmt.Println("Election timeout")
			s.checkElectionResults()
			return
		}
	}
}

func (s *Server) sendMessage(msg string, address string) {
	message := "[" + s.priority + "]" + msg
	fmt.Println("Send msg", message, address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		//println("Connect to server failed:", err.Error())
		return
	}
	fmt.Fprintf(conn, message+"\n")
}

func (s *Server) checkElectionResults() {
	switch {
	case s.state == "election":
		s.setState("king")
		go s.sendToAll("IMTHEKING")
	//case s.state == "waiting-for-king":
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message Received:", string(message))
	defer conn.Close()

	msg, ok  := s.parseProtocol(message)
	if !ok {
		return
	}

	switch {
	case msg.msg == "ALIVE?\n":
		go s.handleAliveMsg(msg)
	case msg.msg == "IMTHEKING\n":
		//go s.handleKingMsg(msg)
		s.setState("worker")
		s.king = msg.priority
		s.pong = true
		fmt.Println("Set king to ", msg.priority)
	case msg.msg == "PING\n":
		if s.state == "king" {
			go s.sendMessage("PONG", msg.Server())
		}
	case msg.msg == "PONG\n":
		if s.state == "worker" && s.king == msg.priority {
			s.pong = true
		}
	case msg.msg == "FINETHANKS\n":
		if s.state == "election" {
			s.setState("waiting-for-king")
			go func() {
				finish := time.After(time.Duration(2*T))
				for {
					select {
					case <-finish:
						fmt.Println("Waiting for king timeout")
						if s.state == "waiting-for-king" {
							//go s.startElections()
							fmt.Println("Election should be started")
						}
					}
				}
			}()
		}
	default:
		return
	}
}

func (s *Server) parseProtocol(msg string) (Message, bool) {
	priority := string(msg[1])
	message := string(msg[3:])
	return Message{
		msg: message,
		priority: priority,
	}, true
}

func (s *Server) handleAliveMsg(msg Message) {
	if s.state == "king" {
		go s.sendMessage("IMTHEKING", msg.Server())
	} else {
		go s.startElections()
		go s.sendMessage("FINETHANKS", msg.Server())
	}
	
}

func (s *Server) sendToAll(msg string) {
	for priority, server := range PRIORITY_MAP {
		if s.priority != priority {
			go s.sendMessage(msg, server)
		}
	}
}

func (s *Server) setState(state string) {
	fmt.Println("Set state:", state)
	s.state = state
}

func (m *Message) Server() string {
	return PRIORITY_MAP[m.priority]
}